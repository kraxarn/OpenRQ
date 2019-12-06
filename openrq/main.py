import sys
from PySide2.QtWidgets import QApplication, QMainWindow, QMenuBar
from PySide2.QtCore import QJsonDocument, QSize

iProject = __import__("project")

class MainWindow(QMainWindow):
	def __init__(self):
		QMainWindow.__init__(self, None)

if __name__ == "__main__":
	app = QApplication()
	window = MainWindow()
	window.resize(1280, 720)
	window.setWindowTitle("OpenRQ")
	window.show()
	p = iProject.Project("OrqProject")
	sys.exit(app.exec_())
