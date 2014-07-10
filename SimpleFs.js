var fs = require("fs")

exports.open = function(file_path){
    return fs.openSync(file_path, "r");
}

exports.read = function(file_path){
    fd = fs.openSync(file_path, "r");
	buffer = new Buffer(fs.statSync(file_path).size);							
	fs.readSync(fd, buffer, 0, buffer.length, 0);
	fs.closeSync(fd);
	return buffer;
}

exports.readAllAsString = function(files){
    var array = [];
    for(var i = 0; i < files.length; i++)
        array.push(exports.read(files[i]).toString());    
    return array;
}

exports.close = function(fd){
    fs.closeSync(fd);		
}

exports.fileExists = function(path){
    return fs.existsSync(path);
}

exports.getFileSize = function(path){
    return fs.statSync(fd.path).size;
}

exports.getFileName = function(path){
	components = path.split("/");
	return components[components.length - 1];					
}

exports.writeFile = function(path, content){
	var buffer = new Buffer(content);
	var f = fs.openSync(path, "w");
	fs.writeSync(f, buffer, 0, buffer.length, 0);
	fs.close(f);
}
