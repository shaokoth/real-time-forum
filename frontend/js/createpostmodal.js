// Show create post form
function showCreatePostForm() {
    const modal = document.getElementById('create-post-form');
    modal.style.display = 'flex';
    populateCategoriesForPost();
}

// Close create post form
function closeCreatePostForm() {
    document.getElementById('create-post-form').style.display = 'none';
}

document.getElementById('create-post-form').addEventListener('click', function(e) {
    if (e.target === this) {
        closeCreatePostForm();
    }
});

