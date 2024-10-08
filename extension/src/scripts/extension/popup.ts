import { saveTabs } from "../api/urls";
import { SYNCTAB_API_URL } from "../utils/config";

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

const displayTabs = () => {
  if (!tabContainer) return;
  tabContainer.innerHTML = "";

  chrome.storage.local.get("user", async function (payload) {
    try {
      const res = await fetch(`${SYNCTAB_API_URL}/urls/url-group`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${payload.user.token}`,
        },
      });

      const data = await res.json();

      const savedTabs =
        (data.length > 0 &&
          data?.sort((a: TabGroup, b: TabGroup) => {
            return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
          })) ||
        [];

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
    } catch (error) {
      console.error("Error:", error);
      const syncTabUrl = chrome.runtime.getURL("synctab.html");
      chrome.tabs.create({ url: "login.html" });
    }
  });
};

document.addEventListener("DOMContentLoaded", displayTabs);

const onClickUrl = async (groupId: string, tabId: number, url: string) => {
  const tabItem = document.getElementById(tabId.toString());
  const tabGroup = document.querySelector(`.tab-group[data-group-id="${groupId}"]`);

  chrome.storage.local.get("user", async function (payload) {
    try {
      const res = await fetch(`${SYNCTAB_API_URL}/urls/url-group/${tabId}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${payload.user.token}`,
        },
      });

      const data = await res.json();
      console.log(data);

      displayTabs();
    } catch (error) {
      console.error("Error:", error);
    }
  });

  chrome.tabs.create({ url: url, active: false });
};

const deleteTab = async (groupId: string, tabId: number) => {
  const tabItem = document.getElementById(tabId.toString());
  const tabGroup = document.querySelector(`.tab-group[data-group-id="${groupId}"]`);

  chrome.storage.local.get("user", async function (payload) {
    try {
      const res = await fetch(`${SYNCTAB_API_URL}/urls/url-group/${tabId}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${payload.user.token}`,
        },
      });

      const data = await res.json();
      console.log(data);

      displayTabs();
    } catch (error) {
      console.error("Error:", error);
    }
  });
};

chrome.runtime.onMessage.addListener((request) => {
  if (request.action === "initiateSync") {
    chrome.tabs.query({}, async function (tabs) {
      const syncTabUrl = chrome.runtime.getURL("synctab.html");
      const findAllSyncTabs = tabs.filter((tab) => tab.url === syncTabUrl);

      if (findAllSyncTabs.length > 1) {
        findAllSyncTabs.forEach((tab) => !tab.active && chrome.tabs.remove(tab.id || 0));
      }

      const tabsToSave = tabs.filter(
        (tab) => !tab.url?.includes("chrome://") && !tab.url?.includes("chrome-extension://")
      );

      if (tabsToSave.length > 0) {
        const res = await saveTabs(tabsToSave, request.data);

        if (res) {
          displayTabs();

          await chrome.tabs.remove(tabsToSave.map((tab) => tab.id || 0)).then(() => {
            const syncTabs = tabs.filter((tab) => tab.url === syncTabUrl);

            if (syncTabs.length > 1) {
              syncTabs.forEach((tab) => !tab?.active && chrome.tabs.remove(tab?.id || 0));
            }
          });
        }
      }
    });
  } else if (request.action === "displayTabs") {
    // displayTabs();
  } else {
    console.log("No action");
  }
});
