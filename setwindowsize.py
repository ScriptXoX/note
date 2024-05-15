import ctypes
import win32gui
import win32api,win32con
import pyautogui
import pydirectinput

# #获取一个窗口的句柄
HWND = win32gui.FindWindow(None,"Guild Wars 2")


# #设置位置和大小 
win32gui.SetWindowPos(HWND, None, 0, 0, 1600, 1000, win32con.SWP_NOSENDCHANGING|win32con.SWP_SHOWWINDOW)

win32gui.SetForegroundWindow(HWND)


