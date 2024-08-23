const SYNCTAB_URL = "synctab.html";

const createParams = {
  url: chrome.runtime.getURL(SYNCTAB_URL),
};

function initiateSync(tabContext: chrome.tabs.Tab) {
  setTimeout(() => {
    chrome.runtime.sendMessage({ action: "initiateSync" });
  }, 500);
}

chrome.action.onClicked.addListener(() => {
  chrome.tabs.create(createParams, (tab) => initiateSync(tab));
});
