function toggleClasses(element, clzes) {
    clzes.forEach(clz => {
        element.classList.toggle(clz);
    });
}
function elementByID(id) {
    return document.getElementById(id);
}
function closestElement(node, selector) {
    return node.closest(selector);
}
