import QtQuick 2.2
import QtQuick.Controls 1.0
import QtQuick.Layouts 1.1
import QtQuick.Dialogs 1.0

ApplicationWindow {
    width: 200
    height: 200
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
        Rectangle {
            Layout.fillHeight: true
            Layout.fillWidth: true
            Flickable {
                clip: true
                boundsBehavior: Flickable.StopAtBounds
                anchors.centerIn: parent
                width: { return Math.min(parent.width, g.width) }
                height: { return Math.min(parent.height, g.height) }
                contentWidth: g.width; contentHeight: g.height
                flickableDirection: Flickable.HorizontalAndVerticalFlick
                Grid {
                    id: g
                    objectName: "grid"
                    spacing: 1
                }
            }
        }
        Rectangle {
            id: botBorder
            height: 1
            Layout.minimumHeight: 1
            Layout.fillWidth: true
            color: "black"
        }
    }
}
