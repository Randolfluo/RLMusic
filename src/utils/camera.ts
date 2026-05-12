import { Camera, CameraResultType, CameraSource } from "@capacitor/camera";

export const capacitorPickImage = async (): Promise<File> => {
  const image = await Camera.getPhoto({
    quality: 90,
    allowEditing: false,
    resultType: CameraResultType.Uri,
    source: CameraSource.Prompt,
  });

  const response = await fetch(image.webPath!);
  const blob = await response.blob();
  const fileName = image.path?.split("/").pop() || `image.${image.format}`;

  return new File([blob], fileName, { type: blob.type });
};
