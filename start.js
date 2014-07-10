//var sfs = require("./SimpleFs.js");
var code_patch = require("./code_patch.js");

function Main(args){
    if(args.length != 6){
        console.log("* Command not Found");
        console.log("* Use node start.js version_one.file developed.file new_update.file output.file");
        return;
    }
    console.log("Start Working");
}

Main(process.argv);
