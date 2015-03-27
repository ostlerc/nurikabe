import QtQuick 2.2
import QtQuick.Controls 1.0
import QtQuick.Layouts 1.1

ApplicationWindow {
    id: app
    width: 300
    height: 300
    color: "white"
    ColumnLayout {
        id: screen
        anchors.fill: parent
        spacing: 0
        Rectangle {
            id: statusRect
            Layout.fillWidth: true
            Layout.preferredHeight: statusRow.height + 10
            border.color: "black"
            border.width: 1
            RowLayout {
                id: statusRow
                anchors { verticalCenter: parent.verticalCenter; margins: 5; horizontalCenter: parent.horizontalCenter }
                Text {
                    id: statusText
                    objectName: "statusText"
                    text: "Nurikabe"
                }
            }
        }

        Loader {
            id: pageLoader
            objectName: "pageLoader"
            Layout.fillHeight: true
            Layout.fillWidth: true
            source: "game.qml"
        }

        Rectangle {
            id: botBorder
            height: statusRow.height + 10
            border.color: "black"
            border.width: 1
            Layout.fillWidth: true
            RowLayout{
                anchors.fill: parent
                Button {
                    id: menuBtn
                    anchors.left: parent.left
                    text: "Menu"
                    onClicked: window.mainMenuPressed()
                }
                Button {
                    id: nextBtn
                    anchors.right: parent.right
                    text: "Next"
                    visible: false
                }
            }
        }
    }
}
