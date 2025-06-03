// Function to fetch comments for a specific post
async function fetchComments(postId) {
    try {
        const response = await fetch(`/posts/${postId}/comments`);
        if (!response.ok) {
            throw new Error('Failed to fetch comments');
        }
        const comments = await response.json();
        return comments;
    } catch (error) {
        console.error('Error fetching comments:', error);
        return [];
    }
}

// Function to display comments for a post
function displayComments(postId, comments) {
    const commentsContainer = document.querySelector(`#comments-${postId}`);
    if (!commentsContainer) return;

    commentsContainer.innerHTML = '';
    
    if (comments.length === 0) {
        commentsContainer.innerHTML = '<p class="no-comments">No comments yet. Be the first to comment!</p>';
        return;
    }

    comments.forEach(comment => {
        const commentElement = document.createElement('div');
        commentElement.className = 'comment';
        commentElement.innerHTML = `
            <div class="comment-header">
                <span class="comment-author">${comment.author}</span>
                <span class="comment-date">${new Date(comment.created_at).toLocaleDateString()}</span>
            </div>
            <div class="comment-content">${comment.content}</div>
        `;
        commentsContainer.appendChild(commentElement);
    });
}

// Function to add a new comment
async function addComment(postId, content) {
    try {
        const response = await fetch(`/posts/${postId}/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ content }),
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Failed to add comment');
        }

        const newComment = await response.json();
        const comments = await fetchComments(postId);
        displayComments(postId, comments);
        return newComment;
    } catch (error) {
        console.error('Error adding comment:', error);
        throw error;
    }
}

// Function to initialize comments for a post
async function initializeComments(postId) {
    const comments = await fetchComments(postId);
    displayComments(postId, comments);
}

// Function to create comment form
function createCommentForm(postId) {
    const form = document.createElement('div');
    form.className = 'comment-form';
    form.innerHTML = `
        <textarea placeholder="Write a comment..." class="comment-input"></textarea>
        <button onclick="submitComment(${postId})" class="comment-submit">Post Comment</button>
    `;
    return form;
}

// Function to submit a comment
async function submitComment(postId) {
    const commentInput = document.querySelector(`#comments-${postId} .comment-input`);
    const content = commentInput.value.trim();
    
    if (!content) {
        alert('Please enter a comment');
        return;
    }

    try {
        await addComment(postId, content);
        commentInput.value = '';
    } catch (error) {
        alert('Failed to post comment. Please try again.');
    }
}

// Export functions for use in other files
window.fetchComments = fetchComments;
window.displayComments = displayComments;
window.addComment = addComment;
window.initializeComments = initializeComments;
window.createCommentForm = createCommentForm;
window.submitComment = submitComment; 