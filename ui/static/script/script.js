document.addEventListener("DOMContentLoaded", function () {
    document.querySelectorAll('input[type="checkbox"]').forEach(checkbox => {
        checkbox.checked = false;
    });
});


document.querySelector(".uil-bars").addEventListener("click", function () {
    let leftElement = document.querySelector(".left");
    leftElement.style.display = leftElement.style.display === "block" ? "none" : "block";
});

document.querySelector(".uil-clipboard-alt").addEventListener("click", function () {
    let rightElement = document.querySelector(".right");
    rightElement.style.display = rightElement.style.display === "block" ? "none" : "block";
});


document.addEventListener("DOMContentLoaded", function () {
    document.querySelectorAll(".uil-comment-alt-dots").forEach(icon => {
        icon.addEventListener("click", function () {
            let commentId = this.id;
            let commentElement = document.getElementById("post-" + commentId);
            console.log("Comment btn clicked and ID = ", commentId)

            console.log("Post", commentId, "clicked")

            if (commentElement) {
                commentElement.style.display = commentElement.style.display === "block" ? "none" : "block";
            }
        });
    });

    document.querySelectorAll(".reply-btn").forEach(btn => {
        btn.addEventListener("click", function () {
            let replyId = this.id;
            let replyElement = document.getElementById("reply-" + replyId);
            console.log("reply btn clicked and ID = ", replyId)

            if (replyElement) {
                replyElement.style.display = replyElement.style.display === "block" ? "none" : "block";
            }
        });
    });

    document.querySelectorAll('input[name="filteredcategories"]').forEach(btn => {
        btn.addEventListener("click", filterByCategories);
    });

});

document.getElementById('loginButton').addEventListener('click',(e)=>{

})

// Handle post creation form submission
document.getElementById('createPostForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const formData = new FormData(e.target);
    const response = await fetch('/post', {
        method: 'POST',
        body: formData
    });

    if (response.ok) {
        console.log('post created')
        alert('Post created successfully!');
        window.location.href = '/';
    } else if (response.status === 401) {
        alert('Failed to create post: user not logged in.');
        window.location.href = '/login';
    } else {
        alert('Failed to create post.');
        // window.location.reload();
    }
});

// Handle like/dislike form submissions
// Like/dislike handler for both posts and comments
// Like/dislike handler for both posts and comments
document.addEventListener('DOMContentLoaded', function () {
    // Find all like/dislike forms
    const likeForms = document.querySelectorAll('form.like-form').forEach(form => {
        form.addEventListener('submit', async function (e) {
            // Prevent the default form submission
            e.preventDefault();

            // Create a new FormData object from the form
            const formData = new FormData(this);

            // Add the clicked button's value explicitly
            // In modern browsers, we can use e.submitter to get the clicked button
            if (e.submitter) {
                formData.set('type', e.submitter.value);
            } else {
                // Fallback for older browsers - try to determine which button was clicked
                const clickedButton = document.activeElement;
                if (clickedButton && clickedButton.getAttribute('value')) {
                    formData.set('type', clickedButton.getAttribute('value'));
                }
            }

            // Log the form data to confirm
            console.log('Form data being sent:');
            for (let pair of formData.entries()) {
                console.log(pair[0] + ': ' + pair[1]);
            }

            const itemId = formData.get('id');
            const itemType = formData.get('item_type');
            const actionType = formData.get('type'); // 'like' or 'dislike'

            // Reference to the buttons that contain the counts
            const likeButton = this.querySelector('button[value="like"]');
            const dislikeButton = this.querySelector('button[value="dislike"]');
            const likeCountElement = likeButton.querySelector('.like');
            const dislikeCountElement = dislikeButton.querySelector('.dislike');

            // Send AJAX request
            fetch('/likes', {
                method: 'POST',
                body: formData
            })
                .then(response => {
                    // First check if the response is JSON
                    const contentType = response.headers.get('content-type');
                    if (contentType && contentType.includes('application/json')) {
                        // It's JSON, so parse it
                        return response.json().then(data => {
                            return { status: response.status, data, isJson: true };
                        });
                    } else {
                        // Not JSON, get text instead
                        return response.text().then(text => {
                            return { status: response.status, text, isJson: false };
                        });
                    }
                })
                .then(result => {
                    console.log('Response:', result);

                    if (result.status === 200 && result.isJson) {
                        // Success with JSON data
                        const data = result.data;
                        // Update the UI with new counts
                        if (data.likes !== undefined) {
                            likeCountElement.textContent = data.likes;
                        }
                        if (data.dislikes !== undefined) {
                            dislikeCountElement.textContent = data.dislikes;
                        }

                        const clickedButton = actionType === 'like' ? likeButton : dislikeButton;
                        clickedButton.classList.add('active-reaction');
                        setTimeout(() => {
                            clickedButton.classList.remove('active-reaction');
                        }, 500);
                    } else if (result.status === 401) {
                        alert('Please log in to like/dislike posts');
                        window.location.href = '/login';
                    } else {
                        // Handle other error cases
                        const errorMessage = result.isJson && result.data.error
                            ? result.data.error
                            : 'Failed to process like/dislike';
                        console.error('Error response:', result.isJson ? result.data : result.text);
                        alert(errorMessage);
                    }
                })
                .catch(error => {
                    console.error('Network or parsing error:', error);
                    alert('An error occurred while processing your request');
                });
        });
    });
});

//filter posts by categories
async function filterByCategories() {
    const checkboxes = document.querySelectorAll('input[name="filteredcategories"]:checked');


    const selectedCategories = Array.from(checkboxes).map(checkbox => checkbox.value);

    console.log("Categories Sent:", selectedCategories, selectedCategories.length);

    try {
        const response = await fetch("/filtered_posts", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ categories: selectedCategories }),
        });

        const data = await response.json();
        console.log("Filtered Data:", data);
        if (selectedCategories.length === 0) {
            data[0] = "all"
        }
        filter(data)
    } catch (error) {
        console.error("Error filtering:", error);
    }
}

function filter(data) {

    document.querySelectorAll(".feed").forEach(feed => {

        if (data[0] === "all") {
            feed.style.display = "block"
        } else {
            if (data.includes(feed.id)) {
                feed.style.display = "block"
            } else {
                feed.style.display = "none"
            }
        }
    })
}





// Handle comment form submissions
document.querySelectorAll('.comment-form').forEach(form => {
    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const formData = new URLSearchParams(new FormData(event.target));

        try {
            const response = await fetch('/add_comment', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: formData
            });

            if (response.ok) {
                alert('Comment created successfully');
                window.location.reload();
            } else if (response.status === 401) {
                alert('Failed to comment: user not logged in.');
                window.location.href = '/login';
            } else {
                alert('Failed to comment.');
            }
        } catch (error) {
            console.error('Request failed:', error);
        }
    });
});

// Handle reply form submissions
document.querySelectorAll('.reply-form').forEach(form => {
    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const formData = new URLSearchParams(new FormData(event.target));

        try {
            const response = await fetch('/add_reply', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
                body: formData
            });

            if (response.ok) {
                alert('Comment created successfully');
                window.location.reload();
            } else if (response.status === 401) {
                alert('Failed to comment: user not logged in.');
                window.location.href = '/login';
            } else {
                alert('Failed to comment.');
            }
        } catch (error) {
            console.error('Request failed:', error);
        }
    });
});

//profile 
document.getElementById('profile-upload').addEventListener('change', function () {
    let fileName = this.files.length ? this.files[0].name : "No file selected";
    document.getElementById("profile_pic").textContent = "Change Photo to " + fileName;
});

document.getElementById("media").addEventListener("change", function () {
    let fileName = this.files.length ? this.files[0].name : "No file selected";
    document.getElementById("file-name").textContent = fileName;
});

