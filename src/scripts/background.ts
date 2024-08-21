chrome.action.onClicked.addListener(() => {
  chrome.tabs.create({
    url: chrome.runtime.getURL("index.html"),
  });
});

chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
  console.log("Service worker received message from sender %s", sender.id, request);
  sendResponse({ message: "Service worker processed the message" });
});
