/**
 * WebSocket 聊天室 JavaScript
 * 負責處理 WebSocket 連線、訊息發送和接收
 */

class ChatRoom {
    constructor() {
        this.ws = null;
        this.user = '';
        this.reconnectAttempts = 0;
        this.maxReconnectAttempts = 5;
        this.reconnectDelay = 1000;
        this.isManualDisconnect = false;
        
        this.initializeElements();
        this.setupEventListeners();
        this.connectWebSocket();
    }

    initializeElements() {
        this.messagesContainer = document.getElementById('messages');
        this.userInput = document.getElementById('user');
        this.contentInput = document.getElementById('content');
        this.sendButton = document.getElementById('send');
        this.connectionStatus = document.querySelector('.connection-status');
    }

    setupEventListeners() {
        // 發送按鈕點擊事件
        this.sendButton.addEventListener('click', () => this.sendMessage());
        
        // Enter 鍵發送訊息
        this.contentInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                this.sendMessage();
            }
        });

        // 使用者名稱變更時記住
        this.userInput.addEventListener('input', (e) => {
            this.user = e.target.value.trim();
            localStorage.setItem('chatUsername', this.user);
        });

        // 頁面卸載時手動斷線
        window.addEventListener('beforeunload', () => {
            this.isManualDisconnect = true;
            if (this.ws) {
                this.ws.close();
            }
        });

        // 載入儲存的使用者名稱
        const savedUsername = localStorage.getItem('chatUsername');
        if (savedUsername) {
            this.userInput.value = savedUsername;
            this.user = savedUsername;
        }
    }

    connectWebSocket() {
        try {
            this.updateConnectionStatus('connecting', '連線中...');
            this.ws = new WebSocket(`ws://${location.host}/ws`);
            
            this.ws.onopen = () => this.onWebSocketOpen();
            this.ws.onmessage = (event) => this.onWebSocketMessage(event);
            this.ws.onclose = (event) => this.onWebSocketClose(event);
            this.ws.onerror = (error) => this.onWebSocketError(error);
            
        } catch (error) {
            console.error('建立 WebSocket 連線時發生錯誤:', error);
            this.updateConnectionStatus('disconnected', '連線失敗');
        }
    }

    onWebSocketOpen() {
        console.log('WebSocket 連線成功');
        this.reconnectAttempts = 0;
        this.updateConnectionStatus('connected', '已連線');
        this.enableInput();
        
        // 發送歡迎訊息
        this.addSystemMessage('已連接到聊天室，開始聊天吧！');
    }

    onWebSocketMessage(event) {
        try {
            const message = JSON.parse(event.data);
            this.displayMessage(message);
        } catch (error) {
            console.error('解析訊息時發生錯誤:', error);
            this.addSystemMessage('收到無效的訊息格式');
        }
    }

    onWebSocketClose(event) {
        console.log('WebSocket 連線關閉:', event.code, event.reason);
        this.updateConnectionStatus('disconnected', '連線中斷');
        this.disableInput();

        if (!this.isManualDisconnect) {
            this.attemptReconnect();
        }
    }

    onWebSocketError(error) {
        console.error('WebSocket 錯誤:', error);
        this.updateConnectionStatus('disconnected', '連線錯誤');
    }

    attemptReconnect() {
        if (this.reconnectAttempts < this.maxReconnectAttempts) {
            this.reconnectAttempts++;
            const delay = this.reconnectDelay * this.reconnectAttempts;
            
            this.updateConnectionStatus('connecting', `重新連線中... (${this.reconnectAttempts}/${this.maxReconnectAttempts})`);
            
            setTimeout(() => {
                this.connectWebSocket();
            }, delay);
        } else {
            this.updateConnectionStatus('disconnected', '連線失敗，請重新整理頁面');
            this.addSystemMessage('無法重新連線，請重新整理頁面');
        }
    }

    sendMessage() {
        if (!this.ws || this.ws.readyState !== WebSocket.OPEN) {
            this.addSystemMessage('連線中斷，無法發送訊息');
            return;
        }

        const user = this.userInput.value.trim();
        const content = this.contentInput.value.trim();

        if (!user || !content) {
            this.addSystemMessage('請輸入使用者名稱和訊息內容');
            return;
        }

        try {
            const message = { user, content };
            this.ws.send(JSON.stringify(message));
            this.contentInput.value = '';
            this.contentInput.focus();
        } catch (error) {
            console.error('發送訊息時發生錯誤:', error);
            this.addSystemMessage('發送訊息失敗');
        }
    }

    displayMessage(message) {
        const messageElement = document.createElement('div');
        messageElement.className = 'message';
        
        const timestamp = new Date().toLocaleTimeString('zh-TW', { 
            hour12: false, 
            hour: '2-digit', 
            minute: '2-digit' 
        });
        
        messageElement.innerHTML = `
            <span class="username">${this.escapeHtml(message.user)}</span>
            <span class="content">${this.escapeHtml(message.content)}</span>
            <span class="timestamp">${timestamp}</span>
        `;
        
        this.messagesContainer.appendChild(messageElement);
        this.scrollToBottom();
    }

    addSystemMessage(text) {
        const messageElement = document.createElement('div');
        messageElement.className = 'message system-message';
        messageElement.innerHTML = `
            <span class="username">系統</span>
            <span class="content">${this.escapeHtml(text)}</span>
            <span class="timestamp">${new Date().toLocaleTimeString('zh-TW', { hour12: false, hour: '2-digit', minute: '2-digit' })}</span>
        `;
        messageElement.style.borderLeftColor = '#2196F3';
        
        this.messagesContainer.appendChild(messageElement);
        this.scrollToBottom();
    }

    scrollToBottom() {
        this.messagesContainer.scrollTop = this.messagesContainer.scrollHeight;
    }

    updateConnectionStatus(status, text) {
        if (this.connectionStatus) {
            this.connectionStatus.className = `connection-status ${status}`;
            this.connectionStatus.innerHTML = `
                <span class="status-indicator"></span>
                ${text}
            `;
        }
    }

    enableInput() {
        this.sendButton.disabled = false;
        this.userInput.disabled = false;
        this.contentInput.disabled = false;
    }

    disableInput() {
        this.sendButton.disabled = true;
        this.userInput.disabled = true;
        this.contentInput.disabled = true;
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
}

// 頁面載入完成後初始化聊天室
document.addEventListener('DOMContentLoaded', () => {
    console.log('初始化 WebSocket 聊天室');
    window.chatRoom = new ChatRoom();
});

// 導出類別供其他模組使用 (如果需要的話)
if (typeof module !== 'undefined' && module.exports) {
    module.exports = { ChatRoom };
}
