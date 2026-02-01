import CryptoJS from "crypto-js";

// 与后端保持一致
const AES_KEY = CryptoJS.enc.Utf8.parse("12345678901234567890123456789012"); // 32 bytes
const AES_IV = CryptoJS.enc.Utf8.parse("1234567890123456");                  // 16 bytes

/**
 * AES加密
 * @param plaintext 明文
 * @returns Base64编码的密文
 */
export const aesEncrypt = (plaintext: string): string => {
  const encrypted = CryptoJS.AES.encrypt(plaintext, AES_KEY, {
    iv: AES_IV,
    mode: CryptoJS.mode.CBC,
    padding: CryptoJS.pad.Pkcs7,
  });
  return encrypted.toString(); // CryptoJS defaults to Base64
};
