const mixedAlgorithms = [
    "Negative",
    "Rescaling",
    "Shift&Rescale",
    "Bit Plane Slicing",
    "Salt&Pepper noise"
]

const arithmeticOperations = [
    "Addition",
    "Substraction",
    "Multiplication",
    "Division"
]

const bitOperations = [
    "Bitwise Not",
    "Bitwise And",
    "Bitwise Or",
    "Bitwise Xor"
]

const convolutions = [
    "Convolution - Averaging",
    "Convolution - Weighted averaging",
    "Convolution - Four Neighbour Laplacian",
    "Convolution - Eight Neighbour Laplacian",
    "Convolution - Four Neighbour Laplacian Enhancement",
    "Convolution - Eight Neighbour Laplacian Enhancement",
    "Convolution - Roberts One",
    "Convolution - Roberts Two",
    "Convolution - Sobel X",
    "Convolution - Sobel Y"
]

const transformations = [
    "Logarithmic Transformation",
    "Power Law",
    "Random LUT",
]

const algorithmLabels = {
    "Negative": "负片",
    "Rescaling": "重新缩放",
    "Shift&Rescale": "移位和重新缩放",
    "Bit Plane Slicing": "位平面切片",
    "Salt&Pepper noise": "椒盐噪声",
    "Addition": "加法",
    "Substraction": "减法",
    "Multiplication": "乘法",
    "Division": "除法",
    "Bitwise Not": "按位取反",
    "Bitwise And": "按位与",
    "Bitwise Or": "按位或",
    "Bitwise Xor": "按位异或",
    "Convolution - Averaging": "卷积 - 平均",
    "Convolution - Weighted averaging": "卷积 - 加权平均",
    "Convolution - Four Neighbour Laplacian": "卷积 - 四邻域拉普拉斯",
    "Convolution - Eight Neighbour Laplacian": "卷积 - 八邻域拉普拉斯",
    "Convolution - Four Neighbour Laplacian Enhancement": "卷积 - 四邻域拉普拉斯增强",
    "Convolution - Eight Neighbour Laplacian Enhancement": "卷积 - 八邻域拉普拉斯增强",
    "Convolution - Roberts One": "卷积 - 罗伯茨一",
    "Convolution - Roberts Two": "卷积 - 罗伯茨二",
    "Convolution - Sobel X": "卷积 - 索贝尔 X",
    "Convolution - Sobel Y": "卷积 - 索贝尔 Y",
    "Logarithmic Transformation": "对数变换",
    "Power Law": "幂律变换",
    "Random LUT": "随机 LUT"
};

function translateAlgorithmOptions() {
    $("option").each(function() {
        const text = $(this).text();
        if (algorithmLabels[text]) {
            $(this).text(algorithmLabels[text]);
        }
    });
}

function buttonsRules(){
    if ($("#originalImage").attr("src") === "") {
        $("#processImage-button").prop("disabled", true);
        $("#saveImage-button").prop("disabled", true);
    }
    else {
        $("#processImage-button").prop("disabled", false);
    }
    if ($("#resultImage").attr("src") !== "") {
        $("#saveImage-button").show(); // 确保按钮显示
        // $("#saveImage-button").prop("disabled", false);
    }
}

function setupSelectInputField(){
    var algorithmSelect = $("#algorithmSelect");

    $.each(mixedAlgorithms, function(index, value) {
        var option = $("<option>").text(value).val(value);
        algorithmSelect.append(option);
    });

    $.each(arithmeticOperations, function(index, value) {
        var option = $("<option>").text(value).val(value);
        algorithmSelect.append(option);
    });

    $.each(bitOperations, function(index, value) {
        var option = $("<option>").text(value).val(value);
        algorithmSelect.append(option);
    });

    $.each(convolutions, function(index, value) {
        var option = $("<option>").text(value).val(value);
        algorithmSelect.append(option);
    });

    $.each(transformations, function(index, value) {
        var option = $("<option>").text(value).val(value);
        algorithmSelect.append(option);
    });
}

$(document).ready(function() {
    var formData = new FormData();
    var algorithmSelected;
    var file;
    var file1;
    var urlApiCall;

    function saveImage() {
        var imageSrc = $("#resultImage").attr("src");
        if (imageSrc) {
            var a = document.createElement('a');
            a.href = imageSrc;
            a.download = 'resultImage.jpg'; // 设置下载的文件名
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        } else {
            alert("No image to save!");
        }
    }


    function algorithmSelection(){
        $("#algorithmSelect").on('change',function(){
            algorithmSelected = $("#algorithmSelect").val();
            if (arithmeticOperations.indexOf(algorithmSelected) !== -1 || bitOperations.indexOf(algorithmSelected) !== -1 && algorithmSelected !== "Bitwise Not"){
                $('#second-image-div').css("display", "contents");
                alert("Upload a second image to use this algorithm!");
            }
            else
                $('#second-image-div').css("display", "none");
        });
    }

    function uploadSecondImage(){
        $("#customFile2").on("change", function() {
            var input = this;
            if (input.files && input.files[0]) {
                file1 = input.files[0];
                if (file1.type.match('image.*')) {
                    var reader = new FileReader();
                    reader.onload = function(e) {
                        var imageDataUrl = e.target.result;
                        $("#originalImage2").attr("src", imageDataUrl);
                    };
                    reader.readAsDataURL(file1);
                } else {
                    alert("Invalid file format. Please select an image file.");
                }
            }
        });
    }

    function uploadFirstImage(){
        $("#customFile").on("change", function() {
            var input = this;
            if (input.files && input.files[0]) {
                file = input.files[0];
                if (file.type.match('image.*')) {
                    var reader = new FileReader();
                    reader.onload = function(e) {
                        var imageDataUrl = e.target.result;
                        $("#originalImage").attr("src", imageDataUrl);
                        buttonsRules();
                    };
                    reader.readAsDataURL(file);
                } else {
                    alert("Invalid file format. Please select an image file.");
                }
            }
        });
    }

    function processMixedAlgorithms(){

        if (algorithmSelected === "Rescaling" || algorithmSelected === "Shift&Rescale"){
            var scalingFactor = prompt("Insert scaling factor")
            if (!$.isNumeric(scalingFactor)){
                alert("Insert a number for the scaling factor!");
                return;
            }
            formData.append("scalingFactor", scalingFactor);
        }

        if (algorithmSelected === "Shift&Rescale"){
            var shiftingValue = prompt("Insert shifting value")
            if (!$.isNumeric(shiftingValue)){
                alert("Insert a number for the shifting value!");
                return;
            }
            formData.append("shiftingValue", shiftingValue);
        }

        if (algorithmSelected === "Bit Plane Slicing"){
            var nBit = prompt("Insert number of bit plane")
            if (!$.isNumeric(nBit)){
                alert("Insert a number!");
                return;
            }
            formData.append("nBit", nBit);
        }
        urlApiCall = 'http://localhost:8080/imageProcessing/process';

    }

    function processArithmeticAlgorithms(){

        formData.append("secondImage", file1);
        urlApiCall = 'http://localhost:8080/imageProcessing/process/arithmeticOperations';

    }

    function processBitOperationsAlgorithms(){
        if (algorithmSelected !== "Bitwise Not")
            formData.append("secondImage", file1);

        urlApiCall = 'http://localhost:8080/imageProcessing/process/bitOperations';
    }

    function processConvolutions(){


        urlApiCall = 'http://localhost:8080/imageProcessing/process/convolution';
    }

    function processTransformations(){


        var parameter = prompt("Insert parameter")
        if (!$.isNumeric(parameter)){
            alert("Insert a number!");
            return;
        }
        formData.append("param", parameter);
        urlApiCall = 'http://localhost:8080/imageProcessing/process/transformations';

    }

    setupSelectInputField();
    translateAlgorithmOptions();
    buttonsRules();
    algorithmSelection();
    uploadFirstImage();
    uploadSecondImage();

    $("#saveImage-button").on("click", function() {
        saveImage();
    });

    $("#processImage-button").on("click", function() {
        formData.append("image", file);
        formData.append("algorithm", algorithmSelected);

        if (mixedAlgorithms.indexOf(algorithmSelected) !== -1)
            processMixedAlgorithms();

        if (arithmeticOperations.indexOf(algorithmSelected) !== -1)
            processArithmeticAlgorithms();

        if (bitOperations.indexOf(algorithmSelected) !== -1)
            processBitOperationsAlgorithms();

        if (convolutions.indexOf(algorithmSelected) !== -1)
            processConvolutions();

        if (transformations.indexOf(algorithmSelected) !== -1)
            processTransformations();


        $.ajax({
            url: urlApiCall,
            type: 'POST',
            data: formData,
            processData: false,
            contentType: false,
            success: function(response) {
                $("#resultImage").attr("src", "data:image/jpeg;base64," + response);
                buttonsRules();
                formData = new FormData();
            },
            error: function(xhr, status, error) {
                alert("Error: ", xhr);
                formData = new FormData();
            }
        });
    });
});
