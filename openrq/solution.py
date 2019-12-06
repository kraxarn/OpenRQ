from PySide2.QtWidgets import QVector, QSqlQuery, QCryptographicHash, QByteArray

class Solution(Item):
	def __init__(self):
		self.description = ""
		self.shown = True
		return
	
	def saveChanges(self):
		# TODO
		pass

	def getChildren(self, child):
		pass

	def addChild(self, child):
		pass

	def removeChild(self, child):
		pass


	def getHash(self, QByteArray):
		hash = QCryptographicHash(QCryptographicHash.Md5)
		hash.addData(QByteArray().append(QByteArray().setNum(id)))
		hash.addData(QByteArray().append(shown))
		hash.addData(self.description.toUtf8())

		return hash.result()