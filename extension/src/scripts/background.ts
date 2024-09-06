const SYNCTAB_URL = "synctab.html";
const SYNCTAB_LOGIN_URL = "login.html";
const SYNCTAB_API_URL = "http://localhost:5000/api/v1";

type TabResponse = {
  tabs: {
    id: number;
    url: string;
    title: string;
  }[];
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

const createParams = {
  url: chrome.runtime.getURL(SYNCTAB_URL),
};

const initiateSync = async (tabContext: chrome.tabs.Tab) => {
  try {
    const res = await fetch(`${SYNCTAB_API_URL}/urls/url-group`, {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${user.token}`,
      },
    });

    const data: TabResponse = await res.json();

    setTimeout(() => {
      chrome.runtime.sendMessage({ action: "initiateSync", data });
    }, 500);
  } catch (error) {
    console.error("Error:", error);
  }
};

chrome.action.onClicked.addListener(() => {
  isUserLoggedIn().then((isLoggedIn) => {
    isLoggedIn
      ? chrome.tabs.create(createParams, initiateSync)
      : chrome.tabs.create({ url: SYNCTAB_LOGIN_URL });
  });
});
