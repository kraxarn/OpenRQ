
from PySide2.QtWidgets import QString, QVector, QByteArray

class Item():
	def __init__(self):
		self.id = 0
		self.uid = 0
        self.version = 0
        self.shown = True
        self.description = QString()
        self.name = ""
        self.contains_color = [Red, Green, Blue, Purple]
        
	def saveChanges(self):
		pass

	def getChildren(self):
		pass
	
	def addChild(self, child):
		pass
	
	def removeChild(self, child):
		pass

	def getHash(self, QByteArray):
		pass