// Function to fetch comments for a specific post
async function fetchComments(postId) {
    try {
        const response = await fetch(`/comments?post_id=${postId}`);
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

    comments.forEach(comment => {
        const commentElement = document.createElement('div');
        commentElement.className = 'comment';
        commentElement.setAttribute('data-comment-id', comment.comment_id);
        commentElement.innerHTML = `
            <div class="comment-header">
                <span class="comment-author">${comment.author}</span>
                <span class="comment-date">${new Date(comment.created_at).toLocaleString()}</span>
            </div>
            <div class="comment-content">${comment.content}</div>
            <div class="comment-actions">
                <button class="like-btn ${comment.UserLiked === 1 ? 'active' : ''}" onclick="handleCommentLike(${comment.comment_id})">
                    <span class="emoji">👍</span>
                    <span class="like-count">${comment.likes}</span>
                </button>
            </div>
        `;
        commentsContainer.appendChild(commentElement);
    });
}

// Function to add a new comment
async function addComment(postId, content) {
    try {
        const response = await fetch(`/comments`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                post_id: postId,
                content: content 
            }),
            credentials: 'include'
        });
        if (!response.ok) {
            throw new Error('Failed to add comment');
        }

        const newComment = await response.json();
        const comments = await fetchComments(postId);
        displayComments(postId, comments);
        
        // Refresh the posts list to update comment count
        if (typeof fetchPosts === 'function') {
            fetchPosts();
        }
        
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
        <textarea id="comment-input-${postId}" placeholder="Write a comment..." class="comment-input"></textarea>
        <button type="button" onclick="submitComment(${postId})" class="post-btn">Post Comment</button>
    `;
    return form;
}

// Function to submit a comment
async function submitComment(postId) {
    const commentInput = document.getElementById(`comment-input-${postId}`);
    if (!commentInput) {
        console.error('Comment input not found');
        return;
    }
    
    const content = commentInput.value.trim();
    
    if (!content) {
        alert('Please enter a comment');
        return;
    }

    try {
        await addComment(postId, content);
        commentInput.value = '';
    } catch (error) {
        console.error('Error submitting comment:', error);
        alert('Failed to post comment. Please try again.');
    }
}

// Function to handle comment like
async function handleCommentLike(commentId) {
    try {
        const response = await fetch('/comments/like', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ comment_id: commentId }),
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Failed to like comment');
        }

        // Get the post ID from the comment's parent container
        const commentElement = document.querySelector(`[data-comment-id="${commentId}"]`);
        if (!commentElement) {
            console.error('Comment element not found');
            return;
        }
        const postId = commentElement.closest('.post-card').dataset.postId;
        
        // Refresh comments to update counts
        const comments = await fetchComments(postId);
        displayComments(postId, comments);
    } catch (error) {
        console.error('Error liking comment:', error);
    }
}

// Function to handle comment dislike
async function handleCommentDislike(commentId) {
    try {
        const response = await fetch('/comments/dislike', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ comment_id: commentId }),
            credentials: 'include'
        });

        if (!response.ok) {
            throw new Error('Failed to dislike comment');
        }

        // Get the post ID from the comment's parent container
        const commentElement = document.querySelector(`[data-comment-id="${commentId}"]`);
        if (!commentElement) {
            console.error('Comment element not found');
            return;
        }
        const postId = commentElement.closest('.post-card').dataset.postId;
        
        // Refresh comments to update counts
        const comments = await fetchComments(postId);
        displayComments(postId, comments);
    } catch (error) {
        console.error('Error disliking comment:', error);
    }
}

// Export functions for use in other files
window.fetchComments = fetchComments;
window.displayComments = displayComments;
window.addComment = addComment;
window.initializeComments = initializeComments;
window.createCommentForm = createCommentForm;
window.submitComment = submitComment;
window.handleCommentLike = handleCommentLike; 
window.handleCommentDislike = handleCommentDislike; 