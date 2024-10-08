import { SYNCTAB_API_URL } from "../utils/config";

document.addEventListener("DOMContentLoaded", function () {
  const loginForm = document.getElementById("login-form");
  const loginError = document.getElementById("login-error");

  if (!loginForm || !loginError) return;

  loginForm?.addEventListener("submit", function (e) {
    e.preventDefault();
    const email = (document?.getElementById("email") as HTMLInputElement)?.value;
    const password = (document?.getElementById("password") as HTMLInputElement)?.value;

    fetch(`${SYNCTAB_API_URL}/auth/login`, {
      method: "POST",

      headers: {
        "Content-Type": "application/json",
      },

      body: JSON.stringify({ email, password }),
    })
      .then((response) => response.json())
      .then((data) => {
        if (data.token) {
          chrome.storage.local.set({ user: data }, function () {
            chrome.storage.local.get("action", function (payload) {
              chrome.tabs.create({ url: chrome.runtime.getURL("synctab.html") });
              chrome.runtime.sendMessage({ action: payload.action }, () => window.close());
            });
          });
        } else {
          loginError.style.display = "block";
        }
      })
      .catch((error) => {
        console.error("Error:", error);
        loginError.style.display = "block";
      });
  });
});
