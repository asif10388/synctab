import { v4 as uuidv4 } from "uuid";

type Tab = {
  id: number;
  url: string;
  title: string;
};

type TabGroup = {
  tabs: Tab[];
  group_id: string;
  created_at: string;
};

const tabContainer = document?.getElementById("tabs-container");

const convertDate = (date: string) => {
  const newDate = new Date(date);
  return newDate.toLocaleDateString("en-US", {
    month: "long",
    day: "numeric",
    year: "numeric",
    hour: "numeric",
    weekday: "long",
    minute: "numeric",
  });
};

const displayTabs = async () => {
  if (!tabContainer) return;
  tabContainer.innerHTML = "";

  chrome.storage.local.get("savedTabs", function (data) {
    const savedTabs =
      data.savedTabs.sort((a: TabGroup, b: TabGroup) => {
        return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
      }) || [];

    savedTabs.forEach(function (tab: TabGroup) {
      const tabHTML = `				
		<div class="tab-group mb-2" data-group-id="${tab.group_id}">	
			<h3 class="text-2xl font-black">${convertDate(tab.created_at)}</h3>
			<ul id="tab-${tab.group_id}" class="tabs"></ul>
		</div>
		`;

      tabContainer.insertAdjacentHTML("beforeend", tabHTML);

      const tabGroup = document.getElementById(`tab-${tab.group_id}`);
      if (!tabGroup) return;

      tab.tabs.forEach((t: Tab) => {
        const tabItem = document.createElement("li");
        tabItem.id = t.id.toString();

        tabItem.innerHTML = `
        	<span id="url-${t.id?.toString() || ""}" class="font-medium">${t.title || "No title"}</span>
        	<button id="delete-${t.id}">
        		<svg width="24px" height="24px" xmlns="http://www.w3.org/2000/svg" class="ionicon" viewBox="0 0 512 512"><path d="M289.94 256l95-95A24 24 0 00351 127l-95 95-95-95a24 24 0 00-34 34l95 95-95 95a24 24 0 1034 34l95-95 95 95a24 24 0 0034-34z" fill="white"/>
        	</button>
		`;

        tabGroup.appendChild(tabItem);

        const deleteButton = document.getElementById(`delete-${t.id}`);
        if (deleteButton) deleteButton.onclick = () => deleteTab(tab.group_id, t.id || 0);

        const tabToUrl = document.getElementById(`url-${t.id?.toString() || ""}`);
        if (tabToUrl) tabToUrl.onclick = () => onClickUrl(tab.group_id, t.id || 0, t.url || "");
      });
    });
  });
};

document.addEventListener("DOMContentLoaded", displayTabs);

const onClickUrl = async (groupId: string, tabId: number, url: string) => {
  const tabItem = document.getElementById(tabId.toString());
  const tabGroup = document.querySelector(`.tab-group[data-group-id="${groupId}"]`);

  chrome.storage.local.get("savedTabs", function (data) {
    const savedTabs = data.savedTabs.reverse() || [];

    const findTabGroup = savedTabs.find((tabGroup: TabGroup) => tabGroup.group_id === groupId);
    if (!findTabGroup) return;

    const newSavedTabs = findTabGroup.tabs.filter((tab: Tab) => tab.id !== tabId);

    if (newSavedTabs.length === 0) {
      chrome.storage.local.set({
        savedTabs: savedTabs.filter((tabGroup: TabGroup) => tabGroup.group_id !== groupId),
      });

      tabGroup?.remove();
      return;
    }

    chrome.storage.local.set({
      savedTabs: savedTabs.map((tabGroup: TabGroup) => {
        if (tabGroup.group_id === groupId) {
          return {
            ...tabGroup,
            tabs: newSavedTabs,
          };
        }
        return tabGroup;
      }),
    });
  });

  tabItem?.remove();

  chrome.tabs.create({ url: url, active: false });
};

const saveTabs = (tabs: chrome.tabs.Tab[], payload: TabGroup) => {
  chrome.storage.local.get("savedTabs", async function (data) {
    const tabsData = tabs.map((tab) => ({ id: tab.id, title: tab.title, url: tab.url }));

    const tabGroup = [
      { tabs: tabsData, group_id: uuidv4(), created_at: new Date().toISOString() },
    ].concat(payload);

    console.log(tabGroup);

    chrome.storage.local.set(
      {
        savedTabs: [...(data.savedTabs || []), ...tabGroup],
      },

      () => displayTabs()
    );
  });
};

const deleteTab = (groupId: string, tabId: number) => {
  const tabItem = document.getElementById(tabId.toString());
  const tabGroup = document.querySelector(`.tab-group[data-group-id="${groupId}"]`);

  chrome.storage.local.get("savedTabs", function (data) {
    const savedTabs = data.savedTabs.reverse() || [];

    const findTabGroup = savedTabs.find((tabGroup: TabGroup) => tabGroup.group_id === groupId);
    if (!findTabGroup) return;

    const newSavedTabs = findTabGroup.tabs.filter((tab: Tab) => tab.id !== tabId);
    console.log(newSavedTabs);

    if (newSavedTabs.length === 0) {
      chrome.storage.local.set({
        savedTabs: savedTabs.filter((tabGroup: TabGroup) => tabGroup.group_id !== groupId),
      });

      tabGroup?.remove();
      return;
    }

    chrome.storage.local.set({
      savedTabs: savedTabs.map((tabGroup: TabGroup) => {
        if (tabGroup.group_id === groupId) {
          return {
            ...tabGroup,
            tabs: newSavedTabs,
          };
        }

        return tabGroup;
      }),
    });
  });

  tabItem?.remove();
};

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  if (request.action === "initiateSync") {
    console.log(request.data);

    chrome.tabs.query({}, async function (tabs) {
      const syncTabUrl = chrome.runtime.getURL("synctab.html");

      const tabsToSave = tabs.filter(
        (tab) => !tab.url?.includes("chrome://") && !tab.url?.includes("chrome-extension://")
      );

      if (tabsToSave.length > 0) {
        saveTabs(tabsToSave, request.data);

        chrome.tabs.remove(tabsToSave.map((tab) => tab.id || 0)).then(() => {
          const syncTabs = tabs.filter((tab) => tab.url === syncTabUrl);

          if (syncTabs.length > 1) {
            syncTabs.forEach((tab) => !tab?.active && chrome.tabs.remove(tab?.id || 0));
          }
        });

        displayTabs();
      }
    });
  }
});
