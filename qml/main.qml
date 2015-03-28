import QtQuick 2.2
import QtQuick.Controls 1.0
import QtQuick.Layouts 1.1
import QtQuick.Dialogs 1.0

Rectangle {
    Layout.fillHeight: true
    Layout.fillWidth: true
    color: "white"
    Grid {
        anchors.fill: parent
        spacing: 1
        Repeater {
            model: ["e", "m", "h"]

            Repeater {
                property string d: modelData
                model: 10
                Button {
                    text: d + " " + (index + 1)
                    onClicked: window.level(text[0] + text.substring(2))
                }
            }
        }
    }
}
