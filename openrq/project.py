# from PySide2.QtCore import QString, QHash, QPair

iDataContext = __import__("datacontext")

class Project:
	def __init__(self, path):
		self.open = True
		if path.endswith(".orq"):
			path += ".orq"
		self.data = iDataContext.DataContext(path)
		return
	def getVersions():
		pass

	def getData():
		pass
				