
//profile 
document.getElementById('profile-upload').addEventListener('submit', function (e) {
    const fileInput = document.getElementById('profile-upload');
    if (!fileInput.files.length === 0) {
        e.preventDefault()
        alert('please select a file first')
    }
});


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
        window.location.href = '/allposts';
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
document.addEventListener('DOMContentLoaded', () => {
    document.querySelectorAll('.form.like-form').forEach(form => {
        form.addEventListener('submit', async (event) => {
            event.preventDefault();

            const formData = new URLSearchParams();

            // Get common form data
            const idInput = event.target.querySelector('input[name="id"]');
            const itemTypeInput = event.target.querySelector('input[name="item_type"]');

            // Sanity checks to ensure required inputs exist
            if (!idInput || !itemTypeInput) {
                console.error('Missing required form inputs');
                return;
            }

            formData.append('id', idInput.value);
            formData.append('item_type', itemTypeInput.value);
            formData.append('type', event.submitter.value);

            try {
                const response = await fetch('/likes', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/x-www-form-urlencoded',
                    },
                    body: formData
                });

                if (response.ok) {
                    // Determine the correct count span for different scenarios
                    let countSpan = null;
                    const itemType = itemTypeInput.value;
                    const itemId = idInput.value;

                    // Try to find count span for posts first
                    countSpan = document.getElementById(`${event.submitter.value}-count-${itemType}-${itemId}`);

                    // If not found, try to find the nearest parent comment's like/dislike span
                    if (!countSpan) {
                        const parentCommentDiv = event.target.closest('.comment, .nested-comment');
                        if (parentCommentDiv) {
                            const nearestCountSpan = parentCommentDiv.querySelector(
                                `[id$="-count-${itemType}-${itemId}"]`
                            );
                            countSpan = nearestCountSpan;
                        }
                    }

                    if (countSpan) {

                        // alert('Like/Dislike processed successfully');
                        window.location.reload();
                    } else {
                        console.warn('No count span found for this like/dislike');
                        window.location.reload();
                    }
                } else if (response.status === 401) {
                    window.location.href = '/login';
                } else {
                    const errorText = await response.text();
                    console.error('Error:', errorText);
                    alert('Failed to process like/dislike');
                }
            } catch (error) {
                console.error('Request failed:', error);
                alert('Failed to process request');
            }
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

// Handle comment form submissions
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

document.getElementById("media").addEventListener("change", function () {
    let fileName = this.files.length ? this.files[0].name : "No file selected";
    document.getElementById("file-name").textContent = fileName;
});

