import QtQuick 2.2
import QtQuick.Controls 1.0
import QtQuick.Layouts 1.1

ApplicationWindow {
    objectName: "mainwindow"
    width: 300
    height: 350
    color: "white"
    ColumnLayout {
        anchors.fill: parent
        spacing: 0

        Rectangle {
            Layout.fillWidth: true
            Layout.preferredHeight: statusText.height + 10
            border.color: "black"
            border.width: 1
            Text {
                id: statusText
                anchors {
                    verticalCenter: parent.verticalCenter
                    horizontalCenter: parent.horizontalCenter
                    margins: 5;
                }
                objectName: "statusText"
                text: "Nurikabe"
            }
            Text {
                property int moves: 0
                anchors {
                    verticalCenter: parent.verticalCenter
                    left: parent.left
                    margins: 5;
                }
                objectName: "movesText"
                text: "steps: " + moves
                onVisibleChanged: {
                    moves = 0
                }
            }
            Text {
                Timer {
                    interval: 200; running: true; repeat: true
                    onTriggered: timerText.seconds = Math.floor((new Date().getTime() - timerText.start.getTime()) / 1000)
                }
                id: timerText
                property date start: new Date()
                property int seconds: 0
                anchors {
                    verticalCenter: parent.verticalCenter
                    right: parent.right
                    margins: 5;
                }
                objectName: "timerText"
                visible: false
                text: "time: " + seconds
                onVisibleChanged: {
                    timerText.start = new Date()
                    timerText.seconds = 0
                }
            }
        }

        Loader {
            objectName: "pageLoader"
            Layout.fillHeight: true
            Layout.fillWidth: true
            source: "game.qml"
        }

        Rectangle {
            height: statusText.height + 10
            border.color: "black"
            border.width: 1
            Layout.fillWidth: true
            RowLayout{
                anchors.fill: parent
                Button {
                    objectName: "toolBtn"
                    anchors.left: parent.left
                    text: "Menu"
                    onClicked: window.toolButtonClicked()
                }

                Text {
                    objectName: "recordText"
                    anchors.right: parent.right
                    property int moves: 0
                    property int seconds: 0
                    visible: false
                }
            }
        }
    }
}
