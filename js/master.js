/**
 * Scrolls the window to a specified element.
 * @param {Event}  event   The event object that triggered the scroll effect (optional)
 * @param {String} selector The element selector to scroll to
 * @param {int}  offset  Height offset in pixels for something like a nav bar
**/
function scrollTo_(event=null, selector, offset=0) {
    if (event) {
        event.preventDefault();
    }
    let element = document.querySelector(selector);
    window.scroll({
        behavior: 'smooth',
        left: 0,
        top: element.offsetTop-offset
    });
}

document.addEventListener("DOMContentLoaded", event => {
    const match = window.location.href.match(/#[a-zA-Z0-9-]*$/);
    if (match) {
        setTimeout(() => {
            scrollTo_(null, match[0], offset=62);
        }, 100)
        history.pushState("", document.title, window.location.pathname + window.location.search);
    }
    
    document.querySelector("#search input").addEventListener("keypress", e => {
        if (e.key === "Enter") {
            window.location.href = `/search?query=${e.target.value}`;
        }
    });
});
