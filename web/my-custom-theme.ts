
import type { CustomThemeConfig } from '@skeletonlabs/tw-plugin';

export const myCustomTheme: CustomThemeConfig = {
    name: 'my-custom-theme',
    properties: {
		// =~= Theme Properties =~=
		"--theme-font-family-base": `system-ui`,
		"--theme-font-family-heading": `system-ui`,
		"--theme-font-color-base": "var(--color-primary-800)",
		"--theme-font-color-dark": "var(--color-primary-300)",
		"--theme-rounded-base": "9999px",
		"--theme-rounded-container": "8px",
		"--theme-border-base": "1px",
		// =~= Theme On-X Colors =~=
		"--on-primary": "0 0 0",
		"--on-secondary": "0 0 0",
		"--on-tertiary": "0 0 0",
		"--on-success": "0 0 0",
		"--on-warning": "0 0 0",
		"--on-error": "0 0 0",
		"--on-surface": "0 0 0",
		// =~= Theme Colors  =~=
		// primary | #613bf7 
		"--color-primary-50": "231 226 254", // #e7e2fe
		"--color-primary-100": "223 216 253", // #dfd8fd
		"--color-primary-200": "216 206 253", // #d8cefd
		"--color-primary-300": "192 177 252", // #c0b1fc
		"--color-primary-400": "144 118 249", // #9076f9
		"--color-primary-500": "97 59 247", // #613bf7
		"--color-primary-600": "87 53 222", // #5735de
		"--color-primary-700": "73 44 185", // #492cb9
		"--color-primary-800": "58 35 148", // #3a2394
		"--color-primary-900": "48 29 121", // #301d79
		// secondary | #9de1ae 
		"--color-secondary-50": "240 251 243", // #f0fbf3
		"--color-secondary-100": "235 249 239", // #ebf9ef
		"--color-secondary-200": "231 248 235", // #e7f8eb
		"--color-secondary-300": "216 243 223", // #d8f3df
		"--color-secondary-400": "186 234 198", // #baeac6
		"--color-secondary-500": "157 225 174", // #9de1ae
		"--color-secondary-600": "141 203 157", // #8dcb9d
		"--color-secondary-700": "118 169 131", // #76a983
		"--color-secondary-800": "94 135 104", // #5e8768
		"--color-secondary-900": "77 110 85", // #4d6e55
		// tertiary | #3fa0a6 
		"--color-tertiary-50": "226 241 242", // #e2f1f2
		"--color-tertiary-100": "217 236 237", // #d9eced
		"--color-tertiary-200": "207 231 233", // #cfe7e9
		"--color-tertiary-300": "178 217 219", // #b2d9db
		"--color-tertiary-400": "121 189 193", // #79bdc1
		"--color-tertiary-500": "63 160 166", // #3fa0a6
		"--color-tertiary-600": "57 144 149", // #399095
		"--color-tertiary-700": "47 120 125", // #2f787d
		"--color-tertiary-800": "38 96 100", // #266064
		"--color-tertiary-900": "31 78 81", // #1f4e51
		// success | #37b3fc 
		"--color-success-50": "225 244 255", // #e1f4ff
		"--color-success-100": "215 240 254", // #d7f0fe
		"--color-success-200": "205 236 254", // #cdecfe
		"--color-success-300": "175 225 254", // #afe1fe
		"--color-success-400": "115 202 253", // #73cafd
		"--color-success-500": "55 179 252", // #37b3fc
		"--color-success-600": "50 161 227", // #32a1e3
		"--color-success-700": "41 134 189", // #2986bd
		"--color-success-800": "33 107 151", // #216b97
		"--color-success-900": "27 88 123", // #1b587b
		// warning | #d209f8 
		"--color-warning-50": "248 218 254", // #f8dafe
		"--color-warning-100": "246 206 254", // #f6cefe
		"--color-warning-200": "244 194 253", // #f4c2fd
		"--color-warning-300": "237 157 252", // #ed9dfc
		"--color-warning-400": "224 83 250", // #e053fa
		"--color-warning-500": "210 9 248", // #d209f8
		"--color-warning-600": "189 8 223", // #bd08df
		"--color-warning-700": "158 7 186", // #9e07ba
		"--color-warning-800": "126 5 149", // #7e0595
		"--color-warning-900": "103 4 122", // #67047a
		// error | #90df16 
		"--color-error-50": "238 250 220", // #eefadc
		"--color-error-100": "233 249 208", // #e9f9d0
		"--color-error-200": "227 247 197", // #e3f7c5
		"--color-error-300": "211 242 162", // #d3f2a2
		"--color-error-400": "177 233 92", // #b1e95c
		"--color-error-500": "144 223 22", // #90df16
		"--color-error-600": "130 201 20", // #82c914
		"--color-error-700": "108 167 17", // #6ca711
		"--color-error-800": "86 134 13", // #56860d
		"--color-error-900": "71 109 11", // #476d0b
		// surface | #46a1ed 
		"--color-surface-50": "227 241 252", // #e3f1fc
		"--color-surface-100": "218 236 251", // #daecfb
		"--color-surface-200": "209 232 251", // #d1e8fb
		"--color-surface-300": "181 217 248", // #b5d9f8
		"--color-surface-400": "126 189 242", // #7ebdf2
		"--color-surface-500": "70 161 237", // #46a1ed
		"--color-surface-600": "63 145 213", // #3f91d5
		"--color-surface-700": "53 121 178", // #3579b2
		"--color-surface-800": "42 97 142", // #2a618e
		"--color-surface-900": "34 79 116", // #224f74
	}
}