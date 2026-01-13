export async function registerServiceWorker() {
  if (typeof window === "undefined") return;
  if (!("serviceWorker" in navigator)) {
    console.log("Service Worker not supported");
    return;
  }

  try {
    const registration = await navigator.serviceWorker.register("/sw.js", {
      scope: "/",
    });

    registration.addEventListener("updatefound", () => {
      const newWorker = registration.installing;
      if (newWorker) {
        newWorker.addEventListener("statechange", () => {
          if (
            newWorker.state === "installed" &&
            navigator.serviceWorker.controller
          ) {
            // New content available, show update prompt
            console.log("New content available, please refresh.");
          }
        });
      }
    });

    console.log("Service Worker registered successfully");
    return registration;
  } catch (error) {
    console.error("Service Worker registration failed:", error);
  }
}
