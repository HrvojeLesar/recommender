function onTagChange(selectObject) {
    const selectedTag = selectObject.options[selectObject.selectedIndex].value;
    if (selectedTag) {
        const searchParams = new URLSearchParams(window.location.search);
        searchParams.set("tag", selectedTag);
        console.log(searchParams.toString());
        window.location.search = searchParams.toString();
    }
}
