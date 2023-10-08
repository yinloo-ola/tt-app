document.body.addEventListener('htmx:beforeOnLoad', function (evt) {
    if (evt.detail.xhr.status === 409) {
        evt.detail.shouldSwap = true;
        evt.detail.isError = true;
    }
});