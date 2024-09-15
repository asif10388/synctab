const SYNCTAB_URL = "synctab.html";
const SYNCTAB_LOGIN_URL = "login.html";
const SYNCTAB_API_URL = "http://localhost:5000/api/v1";

type Tabs = {
  id: number;
  url: string;
  title: string;
};

type TabResponse = {
  tabs: Tabs[];
  group_id: string;
  created_at: string;
};

const user = {
  token: "",
};

const isUserLoggedIn = async () => {
  const result = await chrome.storage.local.get(["user"]);
  user.token = result.user?.token;
  return !!user.token;
};

const redirectUser = async (callback: (tab: chrome.tabs.Tab) => void) => {
  return isUserLoggedIn().then((isLoggedIn) => {
    isLoggedIn
      ? chrome.tabs.create(createParams, callback)
      : chrome.tabs.create({ url: SYNCTAB_LOGIN_URL });
  });
};

const createParams = {
  url: chrome.runtime.getURL(SYNCTAB_URL),
};

const initiateSync = async () => {
  try {
    setTimeout(() => {
      chrome.runtime.sendMessage({ action: "initiateSync" });
    }, 500);
  } catch (error) {
    console.error("Error:", error);
  }
};

const displayTabs = async () => {
  try {
    setTimeout(() => {
      chrome.runtime.sendMessage({ action: "displayTabs" });
    }, 500);
  } catch (error) {
    console.error("Error:", error);
  }
};

chrome.action.onClicked.addListener(() => redirectUser(initiateSync));

chrome.runtime.onInstalled.addListener(() => {
  chrome.contextMenus.create({
    id: "display-synctab",
    title: "Display SyncTab",
    contexts: ["all"],
  });
});

chrome.contextMenus.onClicked.addListener((info, tab) => {
  if (info.menuItemId === "display-synctab") redirectUser(displayTabs);
});
