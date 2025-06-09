let currentCategory = null;

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

    // Add "All Categories" option
    const allCategoriesElement = document.createElement('div');
    allCategoriesElement.className = 'category-card active';
    allCategoriesElement.innerHTML = `
        <h3>All Categories</h3>
    `;
    allCategoriesElement.onclick = () => {
        // Remove active class from all category cards
        document.querySelectorAll('.category-card').forEach(card => card.classList.remove('active'));
        // Add active class to this card
        allCategoriesElement.classList.add('active');
        currentCategory = null;
        fetchPosts();
    };
    categoriesList.appendChild(allCategoriesElement);

    categories.forEach(category => {
        const categoryElement = document.createElement('div');
        categoryElement.className = 'category-card';
        categoryElement.innerHTML = `
            <h3>${category.name}</h3>
        `;
        categoryElement.onclick = () => {
            // Remove active class from all category cards
            document.querySelectorAll('.category-card').forEach(card => card.classList.remove('active'));
            // Add active class to this card
            categoryElement.classList.add('active');
            filterPostsByCategory(category.name);
        };
        categoriesList.appendChild(categoryElement);
    });
}

// Fetch and display posts
async function fetchPosts() {
    try {
        const url = currentCategory ? `/posts?category=${encodeURIComponent(currentCategory)}` : '/posts';
        const response = await fetch(url);
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
        postElement.setAttribute('data-post-id', post.post_id);
        postElement.innerHTML = `
            <h3>${post.title}</h3>
            <p class="post-content">${post.content}</p>
            <div class="post-meta">
                <span class="post-author">By ${post.nickname}</span>
                <span class="post-date">${new Date(post.created_at).toLocaleDateString()}</span>
            </div>
            <div class="post-categories">
                ${post.categories.map(cat => `<span class="category-tag">${cat}</span>`).join('')}
            </div>
            <div class="post-actions">
                <button onclick="likePost(${post.post_id})" class="like-btn">👍 ${post.likes || 0}</button>
                <button onclick="toggleComments(${post.post_id})" class="comment-btn">💬 Comments (${post.comments_count || 0})</button>
            </div>
            <div id="comments-section-${post.post_id}" class="comments-section" style="display: none;">
                <div id="comments-${post.post_id}" class="comments-container"></div>
                <div class="comment-form">
                    <textarea id="comment-input-${post.post_id}" placeholder="Write a comment..." class="comment-input"></textarea>
                    <button type="button" onclick="submitComment(${post.post_id})" class="post-btn">Post Comment</button>
                </div>
            </div>
        `;
        postsList.appendChild(postElement);
    });
}

// Filter posts by category
async function filterPostsByCategory(categoryName) {
    try {
        currentCategory = categoryName;
        const response = await fetch(`/posts?category=${encodeURIComponent(categoryName)}`);
        const posts = await response.json();
         displayPosts(posts);
    } catch (error) {
        console.error('Error filtering posts:', error);
    }
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
            fetchPosts();
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

// Toggle comments section visibility
function toggleComments(postId) {
    const commentsSection = document.getElementById(`comments-section-${postId}`);
    if (!commentsSection) return;

    if (commentsSection.style.display === 'none') {
        commentsSection.style.display = 'block';
        initializeComments(postId);
    } else {
        commentsSection.style.display = 'none';
    }
}

// Initialize the page
document.addEventListener('DOMContentLoaded', () => {
    fetchCategories();
    fetchPosts();
}); 