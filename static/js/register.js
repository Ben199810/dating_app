/**
 * 註冊表單驗證與處理
 * 包含18+年齡驗證、表單驗證、API整合
 */

class RegisterFormHandler {
    constructor() {
        this.form = document.getElementById('register-form');
        this.ageModal = document.getElementById('age-verification-modal');
        this.currentStep = 1;
        this.maxStep = 3;
        this.formData = {};
        
        this.init();
    }

    init() {
        this.bindEvents();
        this.setupValidation();
        this.showAgeVerification();
    }

    bindEvents() {
        // 年齡驗證事件
        document.getElementById('age-confirm-yes')?.addEventListener('click', () => {
            this.confirmAge(true);
        });

        document.getElementById('age-confirm-no')?.addEventListener('click', () => {
            this.confirmAge(false);
        });

        // 表單步驟導航
        document.getElementById('next-step-1')?.addEventListener('click', (e) => {
            e.preventDefault();
            this.nextStep();
        });

        document.getElementById('next-step-2')?.addEventListener('click', (e) => {
            e.preventDefault();
            this.nextStep();
        });

        document.getElementById('prev-step-2')?.addEventListener('click', (e) => {
            e.preventDefault();
            this.prevStep();
        });

        document.getElementById('prev-step-3')?.addEventListener('click', (e) => {
            e.preventDefault();
            this.prevStep();
        });

        // 表單提交
        this.form?.addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleSubmit();
        });

        // 即時驗證
        this.setupRealTimeValidation();
    }

    showAgeVerification() {
        if (this.ageModal) {
            this.ageModal.style.display = 'flex';
        }
    }

    confirmAge(isAdult) {
        if (!isAdult) {
            this.showError('很抱歉，本服務僅限18歲以上成人使用。');
            setTimeout(() => {
                window.location.href = 'https://www.google.com';
            }, 2000);
            return;
        }

        // 記錄年齡驗證
        sessionStorage.setItem('ageVerified', 'true');
        this.hideAgeVerification();
    }

    hideAgeVerification() {
        if (this.ageModal) {
            this.ageModal.style.display = 'none';
        }
    }

    setupValidation() {
        // 自定義驗證規則
        this.validationRules = {
            email: {
                required: true,
                pattern: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
                message: '請輸入有效的電子郵件地址'
            },
            password: {
                required: true,
                minLength: 8,
                pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/,
                message: '密碼必須包含至少8個字符，包括大小寫字母、數字和特殊符號'
            },
            confirmPassword: {
                required: true,
                match: 'password',
                message: '確認密碼與密碼不相符'
            },
            name: {
                required: true,
                minLength: 2,
                maxLength: 20,
                pattern: /^[\u4e00-\u9fa5a-zA-Z\s]+$/,
                message: '姓名僅可包含中文、英文字母和空格，長度2-20字符'
            },
            birthDate: {
                required: true,
                custom: (value) => this.validateAge(value),
                message: '您必須年滿18歲才能註冊'
            },
            gender: {
                required: true,
                message: '請選擇您的性別'
            },
            termsAccepted: {
                required: true,
                message: '請同意服務條款和隱私政策'
            }
        };
    }

    setupRealTimeValidation() {
        // 為所有輸入欄位添加即時驗證
        const inputs = this.form?.querySelectorAll('input, select');
        inputs?.forEach(input => {
            input.addEventListener('blur', () => {
                this.validateField(input.name, input.value);
            });

            input.addEventListener('input', () => {
                // 清除之前的錯誤狀態
                this.clearFieldError(input.name);
                
                // 特殊欄位的即時處理
                if (input.name === 'confirmPassword') {
                    this.validateConfirmPassword();
                }
            });
        });

        // 密碼強度指示器
        const passwordInput = document.getElementById('password');
        passwordInput?.addEventListener('input', () => {
            this.updatePasswordStrength(passwordInput.value);
        });
    }

    validateField(fieldName, value) {
        const rule = this.validationRules[fieldName];
        if (!rule) return true;

        const errors = [];

        // 必填驗證
        if (rule.required && (!value || value.trim() === '')) {
            errors.push('此欄位為必填');
        }

        if (value && value.trim() !== '') {
            // 長度驗證
            if (rule.minLength && value.length < rule.minLength) {
                errors.push(`最少需要 ${rule.minLength} 個字符`);
            }
            if (rule.maxLength && value.length > rule.maxLength) {
                errors.push(`最多允許 ${rule.maxLength} 個字符`);
            }

            // 格式驗證
            if (rule.pattern && !rule.pattern.test(value)) {
                errors.push(rule.message);
            }

            // 自定義驗證
            if (rule.custom && !rule.custom(value)) {
                errors.push(rule.message);
            }

            // 匹配驗證
            if (rule.match) {
                const matchField = document.getElementById(rule.match);
                if (matchField && value !== matchField.value) {
                    errors.push(rule.message);
                }
            }
        }

        if (errors.length > 0) {
            this.showFieldError(fieldName, errors[0]);
            return false;
        } else {
            this.clearFieldError(fieldName);
            return true;
        }
    }

    validateAge(birthDate) {
        if (!birthDate) return false;
        
        const birth = new Date(birthDate);
        const today = new Date();
        const age = today.getFullYear() - birth.getFullYear();
        const monthDiff = today.getMonth() - birth.getMonth();
        
        if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < birth.getDate())) {
            age--;
        }
        
        return age >= 18;
    }

    validateConfirmPassword() {
        const password = document.getElementById('password')?.value;
        const confirmPassword = document.getElementById('confirmPassword')?.value;
        
        if (confirmPassword && password !== confirmPassword) {
            this.showFieldError('confirmPassword', '確認密碼與密碼不相符');
            return false;
        } else {
            this.clearFieldError('confirmPassword');
            return true;
        }
    }

    updatePasswordStrength(password) {
        const strengthMeter = document.querySelector('.password-strength-meter');
        const strengthText = document.querySelector('.password-strength-text');
        
        if (!strengthMeter || !strengthText) return;

        let score = 0;
        let feedback = [];

        // 長度檢查
        if (password.length >= 8) score += 1;
        else feedback.push('至少8個字符');

        // 大寫字母
        if (/[A-Z]/.test(password)) score += 1;
        else feedback.push('包含大寫字母');

        // 小寫字母
        if (/[a-z]/.test(password)) score += 1;
        else feedback.push('包含小寫字母');

        // 數字
        if (/\d/.test(password)) score += 1;
        else feedback.push('包含數字');

        // 特殊字符
        if (/[@$!%*?&]/.test(password)) score += 1;
        else feedback.push('包含特殊符號');

        const strength = ['很弱', '弱', '普通', '強', '很強'][score];
        const colors = ['#e74c3c', '#e67e22', '#f39c12', '#27ae60', '#2ecc71'];
        
        strengthMeter.style.width = `${(score / 5) * 100}%`;
        strengthMeter.style.backgroundColor = colors[score];
        strengthText.textContent = `密碼強度：${strength}`;
        
        if (feedback.length > 0) {
            strengthText.textContent += ` (需要：${feedback.join('、')})`;
        }
    }

    showFieldError(fieldName, message) {
        const field = document.getElementById(fieldName);
        if (!field) return;

        // 添加錯誤樣式
        field.classList.add('error');

        // 顯示錯誤訊息
        let errorElement = field.parentNode.querySelector('.error-message');
        if (!errorElement) {
            errorElement = document.createElement('div');
            errorElement.className = 'error-message';
            field.parentNode.appendChild(errorElement);
        }
        errorElement.textContent = message;
        errorElement.style.display = 'block';
    }

    clearFieldError(fieldName) {
        const field = document.getElementById(fieldName);
        if (!field) return;

        field.classList.remove('error');
        const errorElement = field.parentNode.querySelector('.error-message');
        if (errorElement) {
            errorElement.style.display = 'none';
        }
    }

    nextStep() {
        if (!this.validateCurrentStep()) {
            return;
        }

        if (this.currentStep < this.maxStep) {
            this.hideStep(this.currentStep);
            this.currentStep++;
            this.showStep(this.currentStep);
            this.updateProgressBar();
        }
    }

    prevStep() {
        if (this.currentStep > 1) {
            this.hideStep(this.currentStep);
            this.currentStep--;
            this.showStep(this.currentStep);
            this.updateProgressBar();
        }
    }

    validateCurrentStep() {
        const currentStepElement = document.querySelector(`.form-step[data-step="${this.currentStep}"]`);
        if (!currentStepElement) return true;

        const inputs = currentStepElement.querySelectorAll('input, select');
        let isValid = true;

        inputs.forEach(input => {
            if (!this.validateField(input.name, input.value)) {
                isValid = false;
            }
        });

        return isValid;
    }

    showStep(step) {
        const stepElement = document.querySelector(`.form-step[data-step="${step}"]`);
        if (stepElement) {
            stepElement.classList.add('active');
        }
    }

    hideStep(step) {
        const stepElement = document.querySelector(`.form-step[data-step="${step}"]`);
        if (stepElement) {
            stepElement.classList.remove('active');
        }
    }

    updateProgressBar() {
        const progress = (this.currentStep / this.maxStep) * 100;
        const progressBar = document.querySelector('.progress-fill');
        if (progressBar) {
            progressBar.style.width = `${progress}%`;
        }

        // 更新步驟指示器
        document.querySelectorAll('.step-indicator').forEach((indicator, index) => {
            if (index < this.currentStep) {
                indicator.classList.add('completed');
                indicator.classList.remove('active');
            } else if (index === this.currentStep - 1) {
                indicator.classList.add('active');
                indicator.classList.remove('completed');
            } else {
                indicator.classList.remove('active', 'completed');
            }
        });
    }

    async handleSubmit() {
        // 驗證所有欄位
        if (!this.validateAllFields()) {
            this.showError('請修正表單中的錯誤後再提交');
            return;
        }

        // 收集表單資料
        const formData = this.collectFormData();
        
        try {
            this.showLoading(true);
            
            const response = await this.submitRegistration(formData);
            
            if (response.success) {
                this.showSuccess('註冊成功！正在為您跳轉...');
                
                // 儲存用戶資訊
                localStorage.setItem('user', JSON.stringify(response.data.user));
                localStorage.setItem('token', response.data.token);
                
                // 跳轉到檔案設定頁面
                setTimeout(() => {
                    window.location.href = '/static/html/profile.html';
                }, 2000);
                
            } else {
                throw new Error(response.message || '註冊失敗');
            }
            
        } catch (error) {
            console.error('Registration error:', error);
            this.showError(error.message || '註冊失敗，請稍後再試');
        } finally {
            this.showLoading(false);
        }
    }

    validateAllFields() {
        let isValid = true;
        const inputs = this.form.querySelectorAll('input, select');
        
        inputs.forEach(input => {
            if (!this.validateField(input.name, input.value)) {
                isValid = false;
            }
        });

        return isValid;
    }

    collectFormData() {
        const formData = new FormData(this.form);
        const data = {};
        
        for (let [key, value] of formData.entries()) {
            data[key] = value;
        }

        // 移除確認密碼欄位
        delete data.confirmPassword;
        
        // 處理布爾值
        data.termsAccepted = data.termsAccepted === 'on';
        
        return data;
    }

    async submitRegistration(data) {
        const response = await fetch('/api/auth/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
        });

        if (!response.ok) {
            const errorData = await response.json().catch(() => ({}));
            throw new Error(errorData.message || `HTTP錯誤: ${response.status}`);
        }

        return await response.json();
    }

    showLoading(show) {
        const submitBtn = document.querySelector('.submit-btn');
        const loadingElement = document.querySelector('.loading-spinner');
        
        if (submitBtn) {
            submitBtn.disabled = show;
            submitBtn.textContent = show ? '註冊中...' : '建立帳戶';
        }
        
        if (loadingElement) {
            loadingElement.style.display = show ? 'block' : 'none';
        }
    }

    showError(message) {
        this.showMessage(message, 'error');
    }

    showSuccess(message) {
        this.showMessage(message, 'success');
    }

    showMessage(message, type) {
        // 移除現有訊息
        const existingMessages = document.querySelectorAll('.message-toast');
        existingMessages.forEach(msg => msg.remove());

        // 建立新訊息
        const messageElement = document.createElement('div');
        messageElement.className = `message-toast ${type}`;
        messageElement.innerHTML = `
            <div class="message-content">
                <span class="message-icon">${type === 'error' ? '❌' : '✅'}</span>
                <span class="message-text">${message}</span>
            </div>
            <button class="message-close" onclick="this.parentElement.remove()">×</button>
        `;

        // 加入到頁面
        document.body.appendChild(messageElement);

        // 自動移除
        setTimeout(() => {
            if (messageElement.parentNode) {
                messageElement.remove();
            }
        }, 5000);
    }
}

// 工具函數
function formatPhoneNumber(input) {
    let value = input.value.replace(/\D/g, '');
    if (value.length >= 4 && value.length <= 7) {
        value = value.replace(/(\d{4})(\d+)/, '$1-$2');
    } else if (value.length > 7) {
        value = value.replace(/(\d{4})(\d{3})(\d+)/, '$1-$2-$3');
    }
    input.value = value;
}

function togglePasswordVisibility(fieldId) {
    const field = document.getElementById(fieldId);
    const toggle = field?.parentNode.querySelector('.password-toggle');
    
    if (field && toggle) {
        if (field.type === 'password') {
            field.type = 'text';
            toggle.textContent = '👁️‍🗨️';
        } else {
            field.type = 'password';
            toggle.textContent = '👁️';
        }
    }
}

// 初始化
document.addEventListener('DOMContentLoaded', () => {
    new RegisterFormHandler();
});

// 匯出供其他模組使用
if (typeof module !== 'undefined' && module.exports) {
    module.exports = RegisterFormHandler;
}