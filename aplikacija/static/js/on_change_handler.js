function onTagChange(selectObject) {
    const selectedTag = selectObject.options[selectObject.selectedIndex].value;
    if (selectedTag) {
        const searchParams = new URLSearchParams(window.location.search);
        searchParams.set("tag", selectedTag);
        window.location.search = searchParams.toString();
    }
}

function onSimilarUsersChange() {
    const searchParams = new URLSearchParams(window.location.search);
    if (searchParams.get("similarUsers")) {
        window.location.search = "";
    } else {
        searchParams.set("similarUsers", "1");
        window.location.search = searchParams.toString();
    }
}
