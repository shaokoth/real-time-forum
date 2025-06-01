// Fetch and display categories
async function fetchCategories() {
    try {
        const response = await fetch('/categories');
        const categories = await response.json();
        displayCategories(categories);
    } catch (error) {
        console.error('Error fetching categories:', error);
    }
}

// Display categories in the categories grid
function displayCategories(categories) {
    const categoriesList = document.getElementById('categories-list');
    categoriesList.innerHTML = '';

    categories.forEach(category => {
        const categoryElement = document.createElement('div');
        categoryElement.className = 'category-card';
        categoryElement.innerHTML = `
            <h3>${category.name}</h3>
            <p>${category.description || ''}</p>
        `;
        categoryElement.onclick = () => filterPostsByCategory(category.id);
        categoriesList.appendChild(categoryElement);
    });
}

// Fetch and display posts
async function fetchPosts() {
    try {
        const response = await fetch('/posts');
        const posts = await response.json();
        displayPosts(posts);
    } catch (error) {
        console.error('Error fetching posts:', error);
    }
}

// Display posts in the posts grid
function displayPosts(posts) {
    const postsList = document.getElementById('posts-list');
    postsList.innerHTML = '';

    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.className = 'post-card';
        postElement.innerHTML = `
            <h3>${post.title}</h3>
            <p class="post-content">${post.content}</p>
            <div class="post-meta">
                <span class="post-author">Posted by ${post.nickname}</span>
                <span class="post-date">${new Date(post.created_at).toLocaleDateString()}</span>
            </div>
            <div class="post-categories">
                ${post.categories.map(cat => `<span class="category-tag">${cat}</span>`).join('')}
            </div>
            <div class="post-actions">
                <button onclick="likePost(${post.post_id})" class="like-btn">👍</button>
                <button onclick="dislikePost(${post.post_id})" class="dislike-btn">👎</button>
            </div>
        `;
        postsList.appendChild(postElement);
    });
}

// Filter posts by category
async function filterPostsByCategory(categoryId) {
    try {
        const response = await fetch(`/posts?category=${categoryId}`);
        const posts = await response.json();
        displayPosts(posts);
    } catch (error) {
        console.error('Error filtering posts:', error);
    }
}

// Show create post form
function showCreatePostForm() {
    document.getElementById('create-post-form').style.display = 'block';
    populateCategoriesForPost();
}

// Populate categories for post creation
async function populateCategoriesForPost() {
    try {
        const response = await fetch('/categories');
        const categories = await response.json();
        const categoriesSelection = document.getElementById('categories-selection');
        
        categoriesSelection.innerHTML = categories.map(category => `
            <div class="category-checkbox">
                <input type="checkbox" id="category-${category.id}" name="categories" value="${category.id}">
                <label for="category-${category.id}">${category.name}</label>
            </div>
        `).join('');
    } catch (error) {
        console.error('Error fetching categories for post:', error);
    }
}

// Handle post creation
document.getElementById('newPostForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const formData = new FormData(e.target);
    const selectedCategories = Array.from(formData.getAll('categories')).map(Number);
    
    const postData = {
        title: formData.get('title'),
        content: formData.get('content'),
        categories: selectedCategories
    };

    try {
        const response = await fetch('/posts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(postData)
        });

        if (response.ok) {
            document.getElementById('create-post-form').style.display = 'none';
            e.target.reset();
            fetchPosts(); // Refresh posts list
        } else {
            const error = await response.text();
            document.getElementById('create-post-error').textContent = error;
            document.getElementById('create-post-error').style.display = 'block';
        }
    } catch (error) {
        console.error('Error creating post:', error);
        document.getElementById('create-post-error').textContent = 'Error creating post';
        document.getElementById('create-post-error').style.display = 'block';
    }
});

// Like post
async function likePost(postId) {
    try {
        const response = await fetch('/posts/like', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ post_id: postId, is_like: true })
        });
        
        if (response.ok) {
            fetchPosts(); // Refresh posts to update like count
        }
    } catch (error) {
        console.error('Error liking post:', error);
    }
}

// Dislike post
async function dislikePost(postId) {
    try {
        const response = await fetch('/posts/dislike', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ post_id: postId, is_like: false })
        });
        
        if (response.ok) {
            fetchPosts(); // Refresh posts to update dislike count
        }
    } catch (error) {
        console.error('Error disliking post:', error);
    }
}

// Initialize the page
document.addEventListener('DOMContentLoaded', () => {
    fetchCategories();
    fetchPosts();
}); 