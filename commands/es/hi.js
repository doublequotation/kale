// requirer.js
var fs = require("fs");

var Require = (filename) => {
  const file = fs.readfileSync(filename, "utf-8");
  vm(file, filename);
};
