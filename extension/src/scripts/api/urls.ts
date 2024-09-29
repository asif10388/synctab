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

export const saveTabs = async (tabs: chrome.tabs.Tab[], payload: TabGroup) => {
  const tabsData = tabs.map((tab) => ({
    url: tab.url || "",
    title: tab.title || "",
  }));

  const token = await chrome.storage.local.get("user");

  if (!token.user) return;

  try {
    const res = await fetch(`${SYNCTAB_API_URL}/urls/url-group`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token.user.token}`,
      },

      body: JSON.stringify(tabsData),
    });

    return res.json();
  } catch (error) {
    console.error("Error:", error);
    return false;
  }
};
