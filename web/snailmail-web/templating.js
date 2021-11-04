// Test to see if the browser supports the HTML template element by checking
// for the presence of the template element's content attribute.
if (!'content' in document.createElement('template')) {
    // Find another way to add the rows to the table because
    // the HTML template element is not supported.
    console.error("<template> is not supported by your browser");
}

function LoadContent(templateId, pageData) {
    // Instantiate the table with the existing HTML tbody
    // and the row with the template
    var content = document.querySelector("#content");
    var template = document.querySelector("#" + templateId);

    // Clone the new row and insert it into the table
    var clone = template.content.cloneNode(true);

    if (pageData != null) {
        for (const key in pageData) {
            const node = clone.querySelector("#" + key);
            node.textContent = pageData[key];
        }
    }

    removeAllChildNodes(content);
    content.appendChild(clone);
}

function removeAllChildNodes(parent) {
    while (parent.firstChild) {
        parent.removeChild(parent.firstChild);
    }
}
