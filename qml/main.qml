import QtQuick 2.2
import QtQuick.Layouts 1.1

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
            objectName: "btnGrid"
            spacing: 15
        }
    }
}
