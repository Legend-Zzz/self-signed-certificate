async function viewFile(fileName) {
    try {
        const response = await fetch(`/view/${fileName}`);

        if (response.ok) {
            const fileContent = await response.text();
            const modalContent = document.getElementById('file-content');
            modalContent.innerHTML = fileContent;

            openModal();
        } else {
            alert('Unable to view file');
        }
    } catch (error) {
        console.error('Error:', error);
    }
}

function toggleDateInfoVisibility() {
    const dateInfoElements = document.querySelectorAll('.date-info');

    dateInfoElements.forEach((dateInfoElement) => {
        const currentDisplay = getComputedStyle(dateInfoElement).display;
        dateInfoElement.style.display = currentDisplay === 'none' ? 'inline' : 'none';
    });

    const hideAllDatesButton = document.getElementById('toggleDateInfoVisibility');

    if (hideAllDatesButton) {
        const buttonText = hideAllDatesButton.textContent;
        hideAllDatesButton.textContent = buttonText === 'Show Dates' ? 'Hide Dates' : 'Show Dates';
    }
}

function openModal() {
    const modal = document.getElementById('myModal');
    modal.style.display = 'grid';
}

function closeModal() {
    const modal = document.getElementById('myModal');
    modal.style.display = 'none';
    document.getElementById('copyStatus').textContent = '';
}

function downloadFile(fileName) {
    const link = document.createElement('a');
    link.href = `/download/${fileName}`;
    link.download = fileName;

    link.click();
}

function deleteFile(fileName) {
    const xhr = new XMLHttpRequest();

    xhr.open('POST', `/delete/${fileName}`, true);

    xhr.onreadystatechange = function () {
        if (xhr.readyState === XMLHttpRequest.DONE) {
            if (xhr.status === 200) {
                window.location.href = '/files';
            } else {
                alert('Failed to delete file: ' + xhr.statusText);
            }
        }
    };

    xhr.send();
}

document.addEventListener('DOMContentLoaded', function () {
    const copyButton = document.getElementById('copyButton');
    const copyStatus = document.getElementById('copyStatus');
    const fileContent = document.getElementById('file-content');

    const clipboard = new ClipboardJS(copyButton, {
        text: function () {
            return fileContent.textContent;
        }
    });

    function updateStatus(message) {
        copyStatus.textContent = message;
        setTimeout(function () {
            copyStatus.textContent = '';
        }, 1000);
    }

    clipboard.on('success', function () {
        updateStatus('Success');
    });

    clipboard.on('error', function () {
        updateStatus('Failed');
    });
});
