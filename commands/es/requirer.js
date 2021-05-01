const fs = require("fs");

const Require = (filename) => {
    const file = fs.readfileSync(filename, "utf-8");
    vm(file, filename);
};
