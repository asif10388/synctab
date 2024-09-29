import { SYNCTAB_API_URL, SYNCTAB_LOGIN_URL, SYNCTAB_URL } from "../utils/config";

const createParams = {
  url: chrome.runtime.getURL(SYNCTAB_URL),
};

const isUserLoggedIn = async () => {
  const result = await chrome.storage.local.get(["user"]);

  if (result.user && result.user.token) {
    try {
      const res = await fetch(`${SYNCTAB_API_URL}/urls/validate`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${result.user.token}`,
        },
      });

      const data = await res.json();

      if (data.message === "valid token") {
        return true;
      }
    } catch (error) {
      console.error("Error:", error);
      return false;
    }
  }

  return false;
};

const redirectUser = async (action: string) => {
  chrome.storage.local.set({ action });

  return isUserLoggedIn().then((isLoggedIn) => {
    isLoggedIn
      ? chrome.tabs.create(createParams, () => {
          setTimeout(async () => {
            await chrome.runtime.sendMessage({ action });
          }, 500);
        })
      : chrome.tabs.create({ url: SYNCTAB_LOGIN_URL });
  });
};

chrome.action.onClicked.addListener(() => redirectUser("initiateSync"));

chrome.runtime.onInstalled.addListener(() => {
  chrome.contextMenus.create({
    contexts: ["all"],
    id: "display-synctab",
    title: "Display SyncTab",
  });
});

chrome.contextMenus.onClicked.addListener(
  (info, tab) => info.menuItemId === "display-synctab" && redirectUser("displayTabs")
);

chrome.runtime.onMessage.addListener(async (request) => {
  switch (request.action) {
    case "initiateSync":
      redirectUser("initiateSync");
      break;
    case "displayTabs":
      redirectUser("displayTabs");
      break;
  }
});
