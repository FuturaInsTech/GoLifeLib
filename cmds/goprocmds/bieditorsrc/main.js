const {
    app,
    BrowserWindow,
    ipcMain,
    dialog
} = require("electron");


const fs = require("fs");

const prettier = require("prettier");

function createWindow() {

    const win = new BrowserWindow({

        width: 1600,
        height: 1000,
		
        webPreferences: {

            preload: __dirname + "/preload.js",
            contextIsolation: true
        }
    });


    //
    // When debug needed uncomment and use OpenDevTools.  Uncomment for DEBUG CONSOLE... 
    // win.webContents.openDevTools();   
    //

    win.loadFile("editor.html");
}

app.whenReady().then(createWindow);

ipcMain.handle(
"open-html",
async () => {

    const result =
        await dialog.showOpenDialog({

            filters: [
                {
                    name: "HTML",
                    extensions: ["html","htm"]
                }
            ]
        });

    if (result.canceled)
        return null;

    const file =
        result.filePaths[0];

    return {

        path: file,

        content:
            fs.readFileSync(
                file,
                "utf8"
            )
    };
});

ipcMain.handle(
"save-html",
async(event,data)=>{

    const answer=
    await dialog.showMessageBox({

        type:"question",

        buttons:[
            "Save",
            "Cancel"
        ],

        defaultId:0,

        cancelId:1,

        message:
            "Overwrite existing HTML file?"
    });

    if(answer.response!==0)
        return false;

    let formattedContent =
        data.content;

    try{

        // Final cleanup of complete HTML
        formattedContent =
            await prettier.format(
                data.content,
                {
                    parser:"html",

                    tabWidth:4,

                    useTabs:false,

                    printWidth:120,

                    htmlWhitespaceSensitivity:
                        "ignore"
                }
            );

    }
    catch(err){

        console.error(
            "Prettier format error:",
            err
        );

        // Continue saving original HTML
    }

    fs.writeFileSync(

        data.path,

        formattedContent,

        "utf8"
    );

    return true;
});