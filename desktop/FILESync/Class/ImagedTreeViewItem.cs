using System;
using System.Windows;
using System.Windows.Input;
using System.Windows.Controls;
using System.Windows.Media;

namespace FILESync.Class
{
    /// <summary>
    /// 简单的Item类，选中和不被选中用两种图标表示
    /// 内置一个StackPanel，包括一个TextBlock显示文本，图片
    /// </summary>
    public class ImagedTreeViewItem : TreeViewItem
    {
        TextBlock text;
        Image img;
        ImageSource srcSelected, srcUnselected;

        /// <summary>
        /// Constructor makes stack with image and text
        /// </summary>
        public ImagedTreeViewItem()
        {
            StackPanel stack = new StackPanel();
            stack.Orientation = Orientation.Horizontal;
            Header = stack;

            img = new Image();
            img.VerticalAlignment = VerticalAlignment.Center;
            img.Margin = new Thickness(0, 0, 2, 0);
            stack.Children.Add(img);

            text = new TextBlock();
            text.VerticalAlignment = VerticalAlignment.Center;
            stack.Children.Add(text);
        }

        /// <summary>
        /// Public porperty for text and images
        /// </summary>
        public string Text
        {
            get { return text.Text; }
            set { text.Text = value; }
        }

        public ImageSource SelectedImage
        {
            get { return srcSelected; }
            set
            {
                srcSelected = value;

                if (IsSelected)
                {
                    img.Source = srcSelected;
                }
            }//end of set
        }//end of public Imagesource SelectedItem

        public ImageSource UnselectedImage
        {
            get { return srcUnselected; }
            set
            {
                srcUnselected = value;
                if (!IsSelected)
                {
                    img.Source = srcUnselected;
                }
            }//end of set
        }//end of public ImageSource UnselectedImage

        /// <summary>
        /// Event override to set image
        /// </summary>
        /// <param name="e"></param>
        protected override void OnSelected(RoutedEventArgs e)
        {
            base.OnSelected(e);
            img.Source = srcSelected;
        }

        protected override void OnUnselected(RoutedEventArgs e)
        {
            base.OnUnselected(e);
            img.Source = srcUnselected;
        }
    }
}
