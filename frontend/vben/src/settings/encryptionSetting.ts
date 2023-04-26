import { isDevMode } from '/@/utils/env'

// System default cache time, in seconds
export const DEFAULT_CACHE_TIME = 60 * 60 * 24 * 7

// aes encryption key
export const cacheCipher = {
  key: '_10101010110101@',
  iv: '@01001010110010_',
}

// Whether the system cache is encrypted using aes
export const enableStorageEncryption = !isDevMode()
