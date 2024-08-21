type Tab = {
  id?: number | undefined;
  url?: string | undefined;
  title?: string | undefined;
};

const tabDiv = document?.getElementById("tabs");

const displayTabs = () => {
  if (!tabDiv) return;
  tabDiv.innerHTML = "";

  chrome.storage.local.get("savedTabs", function (data) {
    const savedTabs = data.savedTabs || [];
    console.log("Saved tabs", savedTabs);

    savedTabs.forEach(function (tab: Tab) {
      var li = document.createElement("li");

      if (li) {
        li.id = tab.id?.toString() || "";
        li.textContent = tab.title || "No title";
        li.onclick = () => onClickUrl(tab.id || 0, tab.url || "");

        tabDiv?.appendChild(li);
      }
    });
  });
};

document.addEventListener("DOMContentLoaded", displayTabs);

const onClickUrl = async (tabId: number, url: string) => {
  chrome.storage.local.get("savedTabs", function (data) {
    const savedTabs = data.savedTabs || [];
    const newSavedTabs = savedTabs.filter((tab: Tab) => tab.id !== tabId);

    chrome.storage.local.set({ savedTabs: newSavedTabs });
  });

  const tabItem = document.getElementById(tabId.toString());
  tabItem?.remove();

  chrome.tabs.create({ url: url, active: false });
};

const saveTabs = (tabs: Tab[]) => {
  chrome.storage.local.get("savedTabs", async function (data) {
    const tabsData = tabs.map((tab) => ({ id: tab.id, title: tab.title, url: tab.url }));
    console.log("Tabs to save", tabsData);

    chrome.storage.local.set({ savedTabs: [...(data.savedTabs || []), ...tabsData] }, function () {
      console.log("Tabs saved", [...(data.savedTabs || []), ...tabsData]);
      displayTabs();
    });
  });
};

const deleteTab = (tabId: number) => {
  chrome.storage.local.get("savedTabs", function (data) {
    const savedTabs = data.savedTabs || [];
    const newSavedTabs = savedTabs.filter((tab: Tab) => tab.id !== tabId);

    chrome.storage.local.set({ savedTabs: newSavedTabs });
  });

  const tabItem = document.getElementById(tabId.toString());
  tabItem?.remove();
};

chrome.tabs.query({}, function (tabs) {
  const currentTab = tabs[0];
  const syncTabUrl = chrome.runtime.getURL("index.html");

  if (currentTab.url === syncTabUrl) {
    return;
  }

  const tabsToSave = tabs.filter(
    (tab) =>
      tab.status === "complete" &&
      !tab.url?.includes("chrome://") &&
      !tab.url?.includes("chrome-extension://")
  );

  if (tabsToSave.length > 0) {
    saveTabs(tabsToSave);
    chrome.tabs.remove(tabsToSave.map((tab) => tab.id || 0));

    displayTabs();
  }
});

chrome.runtime
  .sendMessage({ createTab: false })
  .then((res) => console.log("Message sent from popup to background", res))
  .catch((err) => console.error(err));
