(function () {
    "use strict";

    var content = document.getElementById("content");
    var sidebarBtns = document.querySelectorAll(".sidebar-btn");
    var periodTabs = document.getElementById("period-tabs");
    var periodBtns = document.querySelectorAll(".period-tab");
    var statusBar = document.getElementById("status-bar");
    var statusDot = document.getElementById("status-dot");

    var activeTab = "time";
    var activePeriod = "today";

    // --- Sidebar navigation ---
    sidebarBtns.forEach(function (btn) {
        btn.addEventListener("click", function () {
            var name = btn.dataset.tab;
            if (name === activeTab) return;
            activeTab = name;
            sidebarBtns.forEach(function (b) { b.classList.remove("active"); });
            btn.classList.add("active");
            loadTab(name);
        });
    });

    // --- Period tab navigation ---
    periodBtns.forEach(function (btn) {
        btn.addEventListener("click", function () {
            var period = btn.dataset.period;
            if (period === activePeriod) return;
            activePeriod = period;
            periodBtns.forEach(function (b) { b.classList.remove("active"); });
            btn.classList.add("active");
            loadTimePeriod(period);
        });
    });

    function loadTab(name) {
        content.innerHTML = '<div class="loading">loading...</div>';
        if (name === "time") {
            periodTabs.classList.remove("hidden");
        } else {
            periodTabs.classList.add("hidden");
        }
        switch (name) {
            case "time":
                loadTimePeriod(activePeriod);
                break;
            case "projects":
                loadProjectsList();
                break;
            case "languages":
                loadLanguagesList();
                break;
            case "profile":
                window.go.main.App.GetProfile().then(renderProfile).catch(showError);
                break;
            case "settings":
                loadSettings();
                break;
        }
    }

    // --- Time view ---
    function loadTimePeriod(period) {
        activePeriod = period;
        content.innerHTML = '<div class="loading">loading...</div>';

        var promise;
        switch (period) {
            case "today":
                promise = window.go.main.App.GetToday();
                break;
            case "week":
                promise = window.go.main.App.GetWeek();
                break;
            case "month":
                promise = window.go.main.App.GetMonth();
                break;
            case "year":
                promise = window.go.main.App.GetYear(new Date().getFullYear());
                break;
        }

        promise.then(function (data) {
            if (!data || data.total === "0h 00m") {
                content.innerHTML = '<div class="empty">no data</div>';
                return;
            }

            var html = '<div class="summary-header">';
            html += '<span class="summary-total">' + escapeHtml(data.total) + '</span>';
            html += '</div>';

            if (data.projects && data.projects.length > 0) {
                html += '<div class="section-title">projects</div>';
                html += renderItems(data.projects);
            }

            if (data.languages && data.languages.length > 0) {
                html += '<hr class="divider">';
                html += '<div class="section-title">languages</div>';
                html += renderItems(data.languages);
            }

            content.innerHTML = html;
        }).catch(showError);
    }

    // --- Projects view ---
    function loadProjectsList() {
        content.innerHTML = '<div class="loading">loading...</div>';
        window.go.main.App.GetProjectNames().then(function (items) {
            if (!items || items.length === 0) {
                content.innerHTML = '<div class="empty">no projects</div>';
                return;
            }
            var html = '';
            for (var i = 0; i < items.length; i++) {
                html += '<div class="list-item" data-name="' + escapeAttr(items[i].name) + '">';
                html += '<span class="list-item-name">' + escapeHtml(items[i].name) + '</span>';
                html += '<span class="list-item-duration">' + escapeHtml(items[i].duration) + '</span>';
                html += '</div>';
            }
            content.innerHTML = html;
            bindListItems("project");
        }).catch(showError);
    }

    function loadLanguagesList() {
        content.innerHTML = '<div class="loading">loading...</div>';
        window.go.main.App.GetLanguageNames().then(function (items) {
            if (!items || items.length === 0) {
                content.innerHTML = '<div class="empty">no languages</div>';
                return;
            }
            var html = '';
            for (var i = 0; i < items.length; i++) {
                html += '<div class="list-item" data-name="' + escapeAttr(items[i].name) + '">';
                html += '<span class="list-item-name">' + escapeHtml(items[i].name) + '</span>';
                html += '<span class="list-item-duration">' + escapeHtml(items[i].duration) + '</span>';
                html += '</div>';
            }
            content.innerHTML = html;
            bindListItems("language");
        }).catch(showError);
    }

    function bindListItems(type) {
        var items = content.querySelectorAll(".list-item");
        items.forEach(function (item) {
            item.addEventListener("click", function () {
                var name = item.dataset.name;
                if (type === "project") {
                    loadProjectDetail(name);
                } else {
                    loadLanguageDetail(name);
                }
            });
        });
    }

    // --- Detail views ---
    function loadProjectDetail(name) {
        content.innerHTML = '<div class="loading">loading...</div>';
        window.go.main.App.GetProjectDetail(name).then(function (data) {
            renderDetail(data, "languages", loadProjectsList);
        }).catch(showError);
    }

    function loadLanguageDetail(name) {
        content.innerHTML = '<div class="loading">loading...</div>';
        window.go.main.App.GetLanguageDetail(name).then(function (data) {
            renderDetail(data, "projects", loadLanguagesList);
        }).catch(showError);
    }

    function renderDetail(data, itemsLabel, backFn) {
        var html = '<div class="detail-header">';
        html += '<button class="detail-back" id="detail-back-btn">&larr;</button>';
        html += '<span class="detail-name">' + escapeHtml(data.name) + '</span>';
        html += '</div>';

        html += '<div class="detail-cards">';
        html += detailCard("all time", data.allTime);
        html += detailCard("this month", data.thisMonth);
        html += detailCard("this week", data.thisWeek);
        html += '</div>';

        if (data.items && data.items.length > 0) {
            html += '<div class="section-title">' + escapeHtml(itemsLabel) + '</div>';
            html += renderItems(data.items);
        }

        content.innerHTML = html;
        document.getElementById("detail-back-btn").addEventListener("click", backFn);
    }

    function detailCard(label, value) {
        return '<div class="detail-card">' +
            '<div class="detail-card-label">' + escapeHtml(label) + '</div>' +
            '<div class="detail-card-value">' + escapeHtml(value || "0h 00m") + '</div>' +
            '</div>';
    }

    // --- Settings view ---
    function loadSettings() {
        window.go.main.App.GetVersion().then(function (version) {
            var html = '';
            html += '<div class="summary-header"><span class="summary-label">settings</span></div>';
            html += '<div class="settings-item"><button class="settings-link" id="settings-github">GitHub</button></div>';
            html += '<div class="settings-item"><button class="settings-link" id="settings-donate">Donate</button></div>';
            html += '<div class="settings-item"><span class="settings-label">License</span><span class="settings-value">MIT</span></div>';
            html += '<div class="settings-item"><span class="settings-label">Version</span><span class="settings-value">' + escapeHtml(version) + '</span></div>';
            content.innerHTML = html;

            document.getElementById("settings-github").addEventListener("click", function () {
                window.runtime.BrowserOpenURL("https://github.com/arnaudhrt/devtime");
            });
            document.getElementById("settings-donate").addEventListener("click", function () {
                window.runtime.BrowserOpenURL("https://github.com/sponsors/arnaudhrt");
            });
        }).catch(showError);
    }

    // --- Shared renderers ---
    function renderItems(items) {
        var html = "";
        for (var i = 0; i < items.length; i++) {
            var item = items[i];
            var pct = Math.round(item.percent || 0);
            html += '<div class="item-row">';
            html += '  <span class="item-name">' + escapeHtml(item.name) + '</span>';
            html += '  <div class="item-bar"><div class="item-bar-fill" style="width: ' + pct + '%"></div></div>';
            html += '  <span class="item-duration">' + escapeHtml(item.duration) + '</span>';
            html += '  <span class="item-pct">' + pct + '%</span>';
            html += '</div>';
        }
        return html;
    }

    function renderProfile(data) {
        if (!data || !data.totalTime) {
            content.innerHTML = '<div class="empty">no data</div>';
            return;
        }

        var html = '<div class="summary-header">';
        html += '<span class="summary-label">profile</span>';
        html += '</div>';

        html += '<div class="profile-grid">';
        html += profileCard("tracking since", data.trackingSince);
        html += profileCard("total time", data.totalTime);
        html += profileCard("daily average", data.dailyAverage);
        html += profileCard("days tracked", String(data.daysTracked));
        html += '</div>';

        if (data.topProjects && data.topProjects.length > 0) {
            html += '<hr class="divider">';
            html += '<div class="section-title">top projects</div>';
            for (var i = 0; i < data.topProjects.length; i++) {
                var p = data.topProjects[i];
                html += '<div class="profile-row"><span class="profile-label">' + escapeHtml(p.name) + '</span><span class="profile-value">' + escapeHtml(p.duration) + '</span></div>';
            }
        }

        if (data.topLanguages && data.topLanguages.length > 0) {
            html += '<hr class="divider">';
            html += '<div class="section-title">top languages</div>';
            for (var i = 0; i < data.topLanguages.length; i++) {
                var l = data.topLanguages[i];
                html += '<div class="profile-row"><span class="profile-label">' + escapeHtml(l.name) + '</span><span class="profile-value">' + escapeHtml(l.duration) + '</span></div>';
            }
        }

        content.innerHTML = html;
    }

    function profileCard(label, value) {
        return '<div class="profile-card">' +
            '<div class="profile-card-label">' + escapeHtml(label) + '</div>' +
            '<div class="profile-card-value">' + escapeHtml(value) + '</div>' +
            '</div>';
    }

    // --- Status bar ---
    function refreshStatus() {
        window.go.main.App.GetStatus().then(function (data) {
            if (data && data.active) {
                statusDot.className = "status-active";
            } else {
                statusDot.className = "status-inactive";
            }
        }).catch(function () {
            statusDot.className = "status-inactive";
        });
    }

    statusBar.addEventListener("click", function () {
        activeTab = "status";
        sidebarBtns.forEach(function (b) { b.classList.remove("active"); });
        periodTabs.classList.add("hidden");
        content.innerHTML = '<div class="loading">loading...</div>';
        window.go.main.App.GetStatus().then(renderStatus).catch(showError);
    });

    function renderStatus(data) {
        if (!data) {
            content.innerHTML = '<div class="empty">no data</div>';
            return;
        }

        var statusClass = data.active ? "status-active-dot" : "status-inactive-dot";
        var statusText = data.active ? "active" : "not active";

        var html = '<div class="summary-header"><span class="summary-label">status</span></div>';
        html += '<div class="profile-row"><span class="profile-label">status</span><span class="profile-value"><span class="status-indicator ' + statusClass + '"></span>' + statusText + '</span></div>';

        if (data.project) {
            html += '<div class="profile-row"><span class="profile-label">project</span><span class="profile-value">' + escapeHtml(data.project) + '</span></div>';
            html += '<div class="profile-row"><span class="profile-label">language</span><span class="profile-value">' + escapeHtml(data.language) + '</span></div>';
            html += '<div class="profile-row"><span class="profile-label">editor</span><span class="profile-value">' + escapeHtml(data.editor) + '</span></div>';
            html += '<div class="profile-row"><span class="profile-label">session</span><span class="profile-value">' + escapeHtml(data.session) + '</span></div>';
            if (!data.active) {
                html += '<div class="profile-row"><span class="profile-label">ended at</span><span class="profile-value">' + escapeHtml(data.lastEnd) + '</span></div>';
            }
        } else {
            html += '<div class="empty" style="text-align:left;padding:12px 0">no recent sessions</div>';
        }

        content.innerHTML = html;
    }

    // --- Utilities ---
    function showError(err) {
        content.innerHTML = '<div class="empty">error: ' + escapeHtml(String(err)) + '</div>';
    }

    function escapeHtml(str) {
        var div = document.createElement("div");
        div.textContent = str;
        return div.innerHTML;
    }

    function escapeAttr(str) {
        return str.replace(/&/g, "&amp;").replace(/"/g, "&quot;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
    }

    // --- Init ---
    refreshStatus();
    setInterval(refreshStatus, 30000);
    loadTab("time");
})();
