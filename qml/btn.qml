import QtQuick 2.0
import QtQuick.Controls 1.0
import QtQuick.Controls.Styles 1.0

Button {
    onClicked: window.onDifficultyClicked(file)
    property string color: "grey"
    property string file
    width: 150

    style: ButtonStyle {
        label: Text {
            renderType: Text.NativeRendering
            font.pointSize: 20
            color: "black"
            text: control.text

            verticalAlignment: Text.AlignVCenter
            horizontalAlignment: Text.AlignHCenter
            anchors.fill: parent
        }

        background: Component {
            Rectangle {
                border.width: 1
                radius: 5
                gradient: Gradient {
                    GradientStop { position: 0 ; color: control.pressed ? Qt.darker(control.color) : control.color }
                    GradientStop { position: 1 ; color: control.pressed ? Qt.darker(control.color) : control.color }
                }
                anchors.fill: parent
            }
        }
    }

}
