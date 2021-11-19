// Test to see if the browser supports the HTML template element by checking
// for the presence of the template element's content attribute.
if (!'content' in document.createElement('template')) {
    // Find another way to add the rows to the table because
    // the HTML template element is not supported.
    console.error("<template> is not supported by your browser");
}

function LoadContent(templateId, pageData) {
    console.debug("loading content '" + templateId + "', data: " + JSON.stringify(pageData ? pageData : {}));

    // Instantiate the table with the existing HTML tbody
    // and the row with the template
    var content = document.querySelector("#content");
    var template = document.querySelector("#" + templateId);

    // Clone the new row and insert it into the table
    var clone = template.content.cloneNode(true);

    loadPageData(clone, pageData);

    removeAllChildNodes(content);
    content.appendChild(clone);
}

// UpdateContent that is already loaded into the page.
function UpdateContent(pageData) {
    console.debug("updating content, data: " + JSON.stringify(pageData ? pageData : {}));

    loadPageData(document, pageData);
}

function loadPageData(element, pageData) {
    if (!pageData) {
        return;
    }

    for (const key in pageData) {
        const node = element.getElementById(key);
        if (node === null) {
            console.error("#" + key + " is not a valid element");
            continue
        }
        let objectValue = pageData[key];
        node.textContent = objectValue;
    }
}


function removeAllChildNodes(parent) {
    while (parent.firstChild) {
        parent.removeChild(parent.firstChild);
    }
}
