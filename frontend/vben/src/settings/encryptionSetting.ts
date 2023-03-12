import { isDevMode } from '/@/utils/env'

// System default cache time, in seconds
export const DEFAULT_CACHE_TIME = 60 * 60 * 24 * 7

// aes encryption key
export const cacheCipher = {
  key: 'xyY$#$%^&&56^$&Y&YR^fiu|\\{pop{Poi<:6%@#$^CJ&^^gHU',
  iv: 'knG^&e43w@65&90i90*YH-0+{K][;]\\}|kIOUGftyDTRFuy',
}

// Whether the system cache is encrypted using aes
export const enableStorageEncryption = !isDevMode()
