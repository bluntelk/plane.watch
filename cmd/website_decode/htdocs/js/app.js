$(document).ready(function () {
    let packet = $('#packet');
    let inputForm =$('#input-form');
    inputForm.submit(function (event) {
        var packet = $('#packet').val();
        $.get("/decode", {'packet': packet})
            .done(function (data) {
                $('#result').html(data)
            })
            .fail(function () {
                $('#result').html("Failed to get data")
            });
        event.preventDefault();
    });

    let search = new URLSearchParams(window.location.search)
    const requestedDecode = search.get("q")
    if ("" !== requestedDecode) {
        packet.val(requestedDecode)
        inputForm.submit()
    } else {
        packet.val("*A028009F96887B05FFA000413602;");
    }


    // hook the menu items
    $('li a[data-packet]').on("click", function () {
        packet.val($(this).data('packet'));
        inputForm.submit()
    });
});
var examplePackets = {
    21: [],
    17: [],
    18: [],
    20: [],
    16: []
//        28: ['*E1999863859533;']
};
var lastRandomNumber = -1;
function setExamplePacket(df) {
    var length = examplePackets[df].length;
    var id = parseInt(Math.random() * length, 10);
    if (length > 1) {
        while (lastRandomNumber == id) {
            id = parseInt(Math.random() * length, 10);
        }
        lastRandomNumber = id;
    } else {
        id = 0;
    }
    var packet = examplePackets[df][id];
    var packetField = $('#packet');
    packetField.val(packet);
    packetField.submit();
}

var examples = $('#examples');
for (let key in examplePackets) {
    if (examplePackets.hasOwnProperty(key)) {
        var link = $('<a class="button">DF' + key + '</a>');
        link.val("DF " + key);
        link.click(function (event) {
            setExamplePacket(key);
            event.preventDefault()
        });
        examples.append(link);
        examples.append('&nbsp;');

    }
}
