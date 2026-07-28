const { contextBridge, ipcRenderer } = require("electron");

contextBridge.exposeInMainWorld("api", {

    openHtml: () =>
        ipcRenderer.invoke("open-html"),

    saveHtml: (data) =>
        ipcRenderer.invoke("save-html", data),

});

