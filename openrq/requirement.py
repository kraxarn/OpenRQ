from PySide2.QtWidgets import QString, QVector, QSqlQuery,QCryptographicHash, QByteArray

class Requirement(Item):
	def __init__(self):
		self.rationale = QString()
		self.fitCriterion = QString()
		self.shown = True
		return
    
	def saveChanges(self):
		# TODO
		pass

	def getChildren(self, child):
		# TODO
		pass

	def addChild(self, child):
		pass

	def removeChild(self, child):
		pass
		
	def getHash(self, QByteArray):
		hash = QCryptographicHash(QCryptographicHash.Md5)
		hash.addData(QByteArray().append(QByteArray().setNum(id)))
		hash.addData(QByteArray().append(shown))
		hash.addData(self.rationale.toUtf8())
		hash.addData(self.fitCriterion.toUtf8())

		return hash.result()
		
