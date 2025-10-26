// EasyGo Panel JavaScript

document.addEventListener('DOMContentLoaded', function() {
    // Initialize tooltips
    var tooltipTriggerList = [].slice.call(document.querySelectorAll('[data-bs-toggle="tooltip"]'));
    var tooltipList = tooltipTriggerList.map(function (tooltipTriggerEl) {
        return new bootstrap.Tooltip(tooltipTriggerEl);
    });

    // Auto-hide alerts after 5 seconds
    const alerts = document.querySelectorAll('.alert');
    alerts.forEach(function(alert) {
        setTimeout(function() {
            const bsAlert = new bootstrap.Alert(alert);
            bsAlert.close();
        }, 5000);
    });

    // Confirm delete actions
    const deleteButtons = document.querySelectorAll('[data-action="delete"]');
    deleteButtons.forEach(function(button) {
        button.addEventListener('click', function(e) {
            if (!confirm('Are you sure you want to delete this item?')) {
                e.preventDefault();
            }
        });
    });

    // Service status refresh
    const statusElements = document.querySelectorAll('[data-service-status]');
    if (statusElements.length > 0) {
        setInterval(refreshServiceStatus, 30000); // Refresh every 30 seconds
    }

    // Real-time system stats
    const statsElements = document.querySelectorAll('[data-stat]');
    if (statsElements.length > 0) {
        setInterval(refreshSystemStats, 10000); // Refresh every 10 seconds
    }
});

// Function to refresh service status
function refreshServiceStatus() {
    fetch('/api/services/status')
        .then(response => response.json())
        .then(data => {
            data.forEach(service => {
                const element = document.querySelector(`[data-service-status="${service.name}"]`);
                if (element) {
                    updateStatusBadge(element, service.status);
                }
            });
        })
        .catch(error => console.error('Error refreshing service status:', error));
}

// Function to refresh system statistics
function refreshSystemStats() {
    fetch('/api/system/stats')
        .then(response => response.json())
        .then(data => {
            Object.keys(data).forEach(stat => {
                const element = document.querySelector(`[data-stat="${stat}"]`);
                if (element) {
                    element.textContent = data[stat];
                }
            });
        })
        .catch(error => console.error('Error refreshing system stats:', error));
}

// Function to update status badge
function updateStatusBadge(element, status) {
    element.className = 'badge';
    
    switch (status.toLowerCase()) {
        case 'running':
        case 'active':
            element.classList.add('bg-success');
            element.textContent = 'Running';
            break;
        case 'stopped':
        case 'inactive':
            element.classList.add('bg-danger');
            element.textContent = 'Stopped';
            break;
        case 'loading':
        case 'starting':
            element.classList.add('bg-warning');
            element.textContent = 'Starting';
            break;
        default:
            element.classList.add('bg-secondary');
            element.textContent = 'Unknown';
    }
}

// Service control functions
function startService(serviceName) {
    fetch(`/api/services/${serviceName}/start`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            showAlert('success', `${serviceName} started successfully`);
            refreshServiceStatus();
        } else {
            showAlert('danger', `Failed to start ${serviceName}: ${data.message}`);
        }
    })
    .catch(error => {
        showAlert('danger', `Error starting ${serviceName}: ${error.message}`);
    });
}

function stopService(serviceName) {
    if (confirm(`Are you sure you want to stop ${serviceName}?`)) {
        fetch(`/api/services/${serviceName}/stop`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            }
        })
        .then(response => response.json())
        .then(data => {
            if (data.success) {
                showAlert('success', `${serviceName} stopped successfully`);
                refreshServiceStatus();
            } else {
                showAlert('danger', `Failed to stop ${serviceName}: ${data.message}`);
            }
        })
        .catch(error => {
            showAlert('danger', `Error stopping ${serviceName}: ${error.message}`);
        });
    }
}

function restartService(serviceName) {
    fetch(`/api/services/${serviceName}/restart`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.success) {
            showAlert('success', `${serviceName} restarted successfully`);
            refreshServiceStatus();
        } else {
            showAlert('danger', `Failed to restart ${serviceName}: ${data.message}`);
        }
    })
    .catch(error => {
        showAlert('danger', `Error restarting ${serviceName}: ${error.message}`);
    });
}

// Function to show alerts
function showAlert(type, message) {
    const alertContainer = document.querySelector('.alert-container') || document.querySelector('main');
    
    const alert = document.createElement('div');
    alert.className = `alert alert-${type} alert-dismissible fade show mt-3`;
    alert.setAttribute('role', 'alert');
    alert.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    `;
    
    alertContainer.insertBefore(alert, alertContainer.firstChild);
    
    // Auto-hide after 5 seconds
    setTimeout(() => {
        const bsAlert = new bootstrap.Alert(alert);
        bsAlert.close();
    }, 5000);
}

// Form validation helpers
function validateDomainForm(form) {
    const domain = form.querySelector('input[name="domain"]').value;
    const docroot = form.querySelector('input[name="docroot"]').value;
    
    if (!domain || !docroot) {
        showAlert('danger', 'Please fill in all required fields');
        return false;
    }
    
    // Basic domain validation
    const domainRegex = /^[a-zA-Z0-9][a-zA-Z0-9-]{1,61}[a-zA-Z0-9]\.[a-zA-Z]{2,}$/;
    if (!domainRegex.test(domain)) {
        showAlert('danger', 'Please enter a valid domain name');
        return false;
    }
    
    return true;
}

// Copy to clipboard function
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(function() {
        showAlert('success', 'Copied to clipboard');
    }, function(err) {
        showAlert('danger', 'Failed to copy to clipboard');
    });
}

// Theme toggle (if needed in future)
function toggleTheme() {
    const body = document.body;
    const currentTheme = body.getAttribute('data-theme');
    const newTheme = currentTheme === 'dark' ? 'light' : 'dark';
    
    body.setAttribute('data-theme', newTheme);
    localStorage.setItem('theme', newTheme);
}

// Load saved theme
(function() {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme) {
        document.body.setAttribute('data-theme', savedTheme);
    }
})();