/**
 * 即時聊天室 - 首頁 JavaScript
 * 負責處理 WebSocket 連線檢查和頁面互動
 */

// WebSocket 連線狀態檢查
function checkConnection() {
    const statusIndicator = document.querySelector('.status-indicator');
    
    if (!statusIndicator) {
        console.warn('狀態指示器元素未找到');
        return;
    }
    
    try {
        const ws = new WebSocket(`ws://${window.location.host}/ws`);
        
        ws.onopen = function() {
            statusIndicator.style.background = '#4CAF50';
            statusIndicator.title = '服務運行正常';
            console.log('WebSocket 連線成功');
        };
        
        ws.onerror = function(error) {
            statusIndicator.style.background = '#f44336';
            statusIndicator.title = '服務連線異常';
            console.error('WebSocket 連線錯誤:', error);
        };
        
        ws.onclose = function(event) {
            if (event.code !== 1000) {
                statusIndicator.style.background = '#ff9800';
                statusIndicator.title = '連線已中斷';
                console.log('WebSocket 連線已關閉:', event.code, event.reason);
            }
            ws.close();
        };
        
        // 設定連線超時
        setTimeout(() => {
            if (ws.readyState === WebSocket.CONNECTING) {
                ws.close();
                statusIndicator.style.background = '#f44336';
                statusIndicator.title = '連線超時';
                console.warn('WebSocket 連線超時');
            }
        }, 5000);
        
    } catch (error) {
        statusIndicator.style.background = '#f44336';
        statusIndicator.title = '無法建立連線';
        console.error('建立 WebSocket 連線時發生錯誤:', error);
    }
}

// 平滑滾動功能
function initSmoothScrolling() {
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const targetId = this.getAttribute('href');
            const target = document.querySelector(targetId);
            
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth',
                    block: 'start'
                });
            } else {
                console.warn(`目標元素未找到: ${targetId}`);
            }
        });
    });
}

// 初始化頁面功能
function initializePage() {
    // 檢查 WebSocket 連線狀態
    checkConnection();
    
    // 初始化平滑滾動
    initSmoothScrolling();
    
    console.log('首頁初始化完成');
}

// 頁面載入完成後執行初始化
document.addEventListener('DOMContentLoaded', initializePage);

// 導出函數供其他模組使用 (如果需要的話)
if (typeof module !== 'undefined' && module.exports) {
    module.exports = {
        checkConnection,
        initSmoothScrolling,
        initializePage
    };
}
